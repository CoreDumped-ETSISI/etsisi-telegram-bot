package janitor

import (
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func AdminOnlyMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)

		var chatid int64
		var userid int

		// Text messages
		if update.Message != nil {
			// Do nothing in private chats
			if update.Message.Chat.IsPrivate() {
				return next(ctx)
			}
			chatid = update.Message.Chat.ID
			userid = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			// Callbacks
			// Do nothing in private chats
			if update.CallbackQuery.Message.Chat.IsPrivate() {
				return next(ctx)
			}

			chatid = update.CallbackQuery.Message.Chat.ID
			userid = update.CallbackQuery.From.ID
		}

		m, err := bot.GetChatMember(tb.ChatConfigWithUser{
			ChatID: chatid,
			UserID: userid,
		})

		// Deny if something went wrong.
		if err != nil {
			return err
		}

		if !m.IsAdministrator() {
			return nil
		}

		return next(ctx)
	}
}

func ManagedOnlyMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(tb.Update)

		var chatid int64

		// Text messages
		if update.Message != nil {
			// Do nothing in private chats
			if update.Message.Chat.IsPrivate() {
				return next(ctx)
			}
			chatid = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			// Callbacks
			// Do nothing in private chats
			if update.CallbackQuery.Message.Chat.IsPrivate() {
				return next(ctx)
			}

			chatid = update.CallbackQuery.Message.Chat.ID
		}

		state := ctx.Arg("state").(state.T)

		managed, err := isChatManaged(state, chatid)

		if err != nil {
			return err
		}

		// If it's not managed, stop.
		if !managed {
			return nil
		}

		return next(ctx)
	}
}
