package verify

import (
	"encoding/json"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func Cmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)
	bot := update.State.Bot()

	if IsUserVerified(update.Message.From.ID) {
		m := tb.NewMessage(update.Message.Chat.ID, "Ya estás verificado!")
		_, err := bot.Send(m)
		return err
	}

	token, err := startNewVerification(update.Message.From.ID)

	if err != nil {
		return err
	}

	url := buildVerificationURL(token)

	m := tb.NewMessage(update.Message.Chat.ID, "Pulsa este botón para iniciar el proceso de verificación.")
	m.ReplyMarkup = tb.NewInlineKeyboardMarkup(tb.NewInlineKeyboardRow(tb.NewInlineKeyboardButtonURL("Verificar ✅", url)))
	_, err = bot.Send(m)
	return err
}

func StartListening() {
	state := state.G

	pubsub := state.Redis().Subscribe("USER_VERIFIED")

	ch := pubsub.Channel()

	for msg := range ch {
		var data userVerified
		err := json.Unmarshal([]byte(msg.Payload), &data)

		if err != nil {
			log.WithError(err).WithField("message", msg.Payload).Error("Error deserializing redis message")
			continue
		}

		err = verifyUser(data.UserID)

		if err != nil {
			log.WithError(err).WithField("userid", data.UserID).Error("Error verifying user")
			continue
		}

		sendVerifiedMessage(data.UserID)
	}
}

func sendVerifiedMessage(userid int) {
	// TODO: comando /grupos
	m := tb.NewMessage(int64(userid), "Gracias por verificar tu cuenta! Ya puedes entrar en grupos protegidos. Puedes ver que grupos están disponibles usando el comando /grupos")
	_, _ = state.G.Bot().Send(m)
}
