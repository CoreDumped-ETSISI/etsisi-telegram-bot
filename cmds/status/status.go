package status

import (
	"math"
	"strings"
	"time"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func sendStatus(ctx commander.Context, infra bool) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	services, err := getStatus()

	if err != nil {
		msg := tb.NewMessage(update.Message.Chat.ID, "La API de estado está caída... 🤦‍♂️")
		bot.Send(msg)

		return err
	}

	var sb strings.Builder

	for _, status := range services {
		if infra != status.Infra {
			continue
		}

		if status.Up {
			sb.WriteString("💚")
		} else {
			sb.WriteString("💔")
		}

		sb.WriteString(" ")
		sb.WriteString(status.Name)

		sb.WriteString(" (hace ")

		ago := time.Now().Sub(status.LastCheck)
		human := time.Duration(math.Ceil(ago.Seconds())) * time.Second

		sb.WriteString(human.String())

		sb.WriteString(")\n")
	}

	msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
	_, err = bot.Send(msg)

	return err
}

func StatusCmd(ctx commander.Context) error {
	return sendStatus(ctx, false)
}

func BotStatusCmd(ctx commander.Context) error {
	return sendStatus(ctx, true)
}
