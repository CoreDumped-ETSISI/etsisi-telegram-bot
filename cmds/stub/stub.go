package stub

import (
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func Cmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	msg := tb.NewMessage(update.Message.Chat.ID, "Este comando ha sido temporalmente desactivado.")
	bot.Send(msg)

	return nil
}
