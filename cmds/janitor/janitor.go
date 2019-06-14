package janitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/verify"

	"github.com/globalsign/mgo/bson"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

var (
	ErrBotNotAdmin = errors.New("El bot debe ser administrador de este grupo")
)

func OnMessage(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)

	if update.Message.NewChatMembers != nil {
		for _, m := range *update.Message.NewChatMembers {
			onChatJoin(ctx, update.Message.Chat.ID, m)
		}
	}

	return nil
}

func onChatJoin(ctx commander.Context, chat int64, member tb.User) error {
	update := ctx.Arg("update").(tb.Update)
	state := ctx.Arg("state").(state.T)
	bot := ctx.Arg("bot").(*tb.BotAPI)

	if man, err := isChatManaged(state, chat); !man {
		return err
	}

	if !verify.IsUserVerified(state, member.ID) {
		m := tb.NewMessage(chat, fmt.Sprintf("%v, para acceder a este chat, tienes que verificar tu cuenta de telegram.", member.String()))
		btn := tb.NewInlineKeyboardButtonURL("Verificar ✔", "https://t.me/"+bot.Self.UserName+"?start=verifyme")

		m.ReplyMarkup = tb.NewInlineKeyboardMarkup([]tb.InlineKeyboardButton{btn})

		_, _ = bot.Send(m)

		_, err := bot.KickChatMember(tb.KickChatMemberConfig{
			ChatMemberConfig: tb.ChatMemberConfig{
				ChatID: chat,
				UserID: member.ID,
			},
			UntilDate: update.Message.Time().Add(1 * time.Minute).Unix(),
		})

		return err
	}

	return nil
}

// Ban is a global ban across all chats where the administrator is admin.
func Ban(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	if update.Message.ReplyToMessage == nil {
		m := tb.NewMessage(update.Message.Chat.ID, "Tienes que citar un mensaje de la persona a la que quieres banear.")
		m.ReplyToMessageID = update.Message.MessageID
		_, err := bot.Send(m)
		return err
	}

	who := update.Message.ReplyToMessage.From

	resp, err := bot.KickChatMember(tb.KickChatMemberConfig{
		ChatMemberConfig: tb.ChatMemberConfig{
			ChatID: update.Message.Chat.ID,
			UserID: who.ID,
		},
	})

	if err != nil || !resp.Ok {
		return err
	}

	m := tb.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("El usuario %v ha sido baneado en todos los grupos que son administrados por tí. Para deshacer esto, use el comando /unban_%v",
			who.String(),
			who.ID,
		))
	_, _ = bot.Send(m)

	ev := banEvent{
		UserID:  who.ID,
		ChatID:  update.Message.Chat.ID,
		AdminID: update.Message.From.ID,
	}

	go propagateBan(bot, state, ev, true)

	j, _ := json.Marshal(ev)

	state.Redis().Publish("USER_BANNED", j)

	return nil
}

// Unban unbands someone across all chats where the sender is admin
func Unban(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	whoid := ctx.ArgInt("user")

	resp, err := bot.UnbanChatMember(tb.ChatMemberConfig{
		ChatID: update.Message.Chat.ID,
		UserID: whoid,
	})

	if err != nil || !resp.Ok {
		return err
	}

	m := tb.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("El usuario %v ha sido desbaneado en todos los grupos que son administrados por tí.",
			whoid,
		))
	_, _ = bot.Send(m)

	ev := banEvent{
		UserID:  whoid,
		ChatID:  update.Message.Chat.ID,
		AdminID: update.Message.From.ID,
	}

	go propagateBan(bot, state, ev, false)

	j, _ := json.Marshal(ev)

	state.Redis().Publish("USER_UNBANNED", j)

	return nil
}

func Manage(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	managed, err := isChatManaged(state, update.Message.Chat.ID)

	if err != nil {
		return err
	}

	if !managed {
		err = manageChannel(bot, state, update.Message.Chat.ID, false)

		if err != nil {
			if err == ErrBotNotAdmin {
				m := tb.NewMessage(update.Message.Chat.ID, "El bot debe ser un administrador de este grupo.")
				_, err = bot.Send(m)
				return err
			} else {
				return err
			}
		}

		m := tb.NewMessage(update.Message.Chat.ID, "Este grupo ahora es moderado por mí.")
		_, err = bot.Send(m)
		return err
	} else {
		return showManagementMenu(ctx)
	}

}

func showManagementMenu(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	chat := update.Message.Chat

	sesh := state.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("managed-channels")

	var chanman channelManagement
	err := col.FindId(chat.ID).One(&chanman)

	if err != nil {
		return err
	}

	m := tb.NewMessage(chat.ID, "Administración del canal")

	var emoji rune

	if chanman.Public {
		emoji = '✔'
	} else {
		emoji = '❌'
	}

	buttons := tb.NewInlineKeyboardMarkup([]tb.InlineKeyboardButton{
		tb.NewInlineKeyboardButtonData("Actualizar", fmt.Sprintf("/jannyrefresh %v %v", chat.ID, chanman.Public)),
		tb.NewInlineKeyboardButtonData(fmt.Sprintf("Público %v", emoji), fmt.Sprintf("/jannypublictoggle %v %v", chat.ID, !chanman.Public)),
		tb.NewInlineKeyboardButtonData("Desactivar", fmt.Sprintf("/jannydisable %v", chat.ID)),
	})

	m.ReplyMarkup = buttons

	_, err = bot.Send(m)

	return err
}

func manageChannel(bot *tb.BotAPI, state state.T, chat int64, public bool) error {
	me, err := bot.GetChatMember(tb.ChatConfigWithUser{
		ChatID: chat,
		UserID: bot.Self.ID,
	})

	if err != nil {
		return err
	}

	if !me.IsAdministrator() {
		return ErrBotNotAdmin
	}

	admins, err := bot.GetChatAdministrators(tb.ChatConfig{
		ChatID: chat,
	})

	if err != nil {
		return err
	}

	var ids []int

	for i := range admins {
		ids = append(ids, admins[i].User.ID)
	}

	chatdata, err := bot.GetChat(tb.ChatConfig{
		ChatID: chat,
	})

	if err != nil {
		return err
	}

	// Refresh the invite link.
	_, _ = bot.GetInviteLink(tb.ChatConfig{ChatID: chat})

	man := channelManagement{
		ChatID:   chat,
		AdminsID: ids,
		Public:   public,
		Name:     chatdata.Title,
	}

	sesh := state.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("managed-channels")

	return col.Insert(man)
}

func unmanageChannel(bot *tb.BotAPI, state state.T, chat int64) error {
	sesh := state.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("managed-channels")

	return col.RemoveId(chat)
}

// If status is true, ban the user. Otherwise unban him.
func propagateBan(bot *tb.BotAPI, state state.T, ban banEvent, status bool) error {
	sesh := state.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("managed-channels")

	var groups []channelManagement

	err := col.Find(bson.M{
		"admins_id": ban.AdminID,
	}).All(&groups)

	if err != nil {
		return err
	}

	for i := range groups {
		chatid := groups[i].ChatID

		if chatid == ban.ChatID {
			continue
		}

		if status {
			_, _ = bot.KickChatMember(tb.KickChatMemberConfig{
				ChatMemberConfig: tb.ChatMemberConfig{
					ChatID: chatid,
					UserID: ban.UserID,
				},
			})
		} else {
			_, _ = bot.UnbanChatMember(tb.ChatMemberConfig{
				ChatID: chatid,
				UserID: ban.UserID,
			})
		}
	}

	return nil
}

// Callbacks
// /jannyrefresh {chatid} {public}
func RefreshCb(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	chatraw := ctx.ArgString("chatid")
	chatid, _ := strconv.ParseInt(chatraw, 10, 64)

	pubraw := ctx.ArgString("public")
	pub, _ := strconv.ParseBool(pubraw)

	unmanageChannel(bot, state, chatid)
	manageChannel(bot, state, chatid, pub)

	_, err := bot.AnswerCallbackQuery(tb.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "Base de datos refrescada!",
		ShowAlert:       false,
	})

	return err
}

// /janmnypublictoggle {chatid} {public}
func TogglePublicCb(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	// state := ctx.Arg("state").(state.T)

	chatraw := ctx.ArgString("chatid")
	chatid, _ := strconv.ParseInt(chatraw, 10, 64)

	pubraw := ctx.ArgString("public")
	public, _ := strconv.ParseBool(pubraw)

	var emoji rune

	if public {
		emoji = '✔'
	} else {
		emoji = '❌'
	}

	buttons := tb.NewInlineKeyboardMarkup([]tb.InlineKeyboardButton{
		tb.NewInlineKeyboardButtonData("Actualizar", fmt.Sprintf("/jannyrefresh %v %v", chatid, public)),
		tb.NewInlineKeyboardButtonData(fmt.Sprintf("Público %v", emoji), fmt.Sprintf("/jannypublictoggle %v %v", chatid, !public)),
		tb.NewInlineKeyboardButtonData("Desactivar", fmt.Sprintf("/jannydisable %v", chatid)),
	})

	m := tb.NewEditMessageReplyMarkup(chatid, update.CallbackQuery.Message.MessageID, buttons)

	_, err := bot.Send(m)

	return err
}

// /jannydisable chatid
func DisableCb(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	chatraw := ctx.ArgString("chatid")
	chatid, _ := strconv.ParseInt(chatraw, 10, 64)

	err := unmanageChannel(bot, state, chatid)

	if err == nil {
		m := tb.NewDeleteMessage(chatid, update.CallbackQuery.Message.MessageID)
		_, _ = bot.Send(m)

		_, err := bot.AnswerCallbackQuery(tb.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Este chat ya no es administrado.",
			ShowAlert:       false,
		})

		return err
	}

	_, _ = bot.AnswerCallbackQuery(tb.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	return err
}
