package verify

import (
	"encoding/json"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func Cmd(ctx commander.Context) error {
	update := ctx.Arg("update").(tb.Update)
	bot := ctx.Arg("bot").(*tb.BotAPI)
	state := ctx.Arg("state").(state.T)

	if IsUserVerified(state, update.Message.From.ID) {
		m := tb.NewMessage(update.Message.Chat.ID, "Ya estás verificado!")
		_, err := bot.Send(m)
		return err
	}

	token, err := startNewVerification(state, update.Message.From.ID)

	if err != nil {
		return err
	}

	url := buildVerificationURL(token)

	m := tb.NewMessage(update.Message.Chat.ID, "Pulsa este botón para iniciar el proceso de verificación.")
	m.ReplyMarkup = tb.NewInlineKeyboardMarkup(tb.NewInlineKeyboardRow(tb.NewInlineKeyboardButtonURL("Verificar ✅", url)))
	_, err = bot.Send(m)
	return err
}

func StartListening(state state.T) {
	pubsub := state.Redis().Subscribe("USER_VERIFIED")

	ch := pubsub.Channel()

	for msg := range ch {
		var data userVerified
		err := json.Unmarshal([]byte(msg.Payload), &data)

		if err != nil {
			log.WithError(err).WithField("message", msg.Payload).Error("Error deserializing redis message")
			continue
		}

		err = verifyUser(state, data.UserID)

		if err != nil {
			log.WithError(err).WithField("userid", data.UserID).Error("Error verifying user")
			continue
		}

		sendVerifiedMessage(state, data.UserID)
	}
}

func sendVerifiedMessage(state state.T, userid int) {
	// TODO: comando /grupos
	m := tb.NewMessage(int64(userid), "Gracias por verificar tu cuenta! Ya puedes entrar en grupos protegidos. Puedes ver que grupos están disponibles usando el comando /grupos")
	_, _ = state.Bot().Send(m)
}
