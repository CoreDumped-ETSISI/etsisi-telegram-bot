package stub

import (
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func Cmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)

	msg := tb.NewMessage(update.Message.Chat.ID, "Este comando ha sido temporalmente desactivado.")
	update.State.Bot().Send(msg)

	return nil
}
