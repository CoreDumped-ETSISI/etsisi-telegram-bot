package news

import (
	"strings"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func handleCmd(ctx commander.Context, feed string, limit int) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	news, err := fetchFeed(feed)

	if err != nil {
		msg := tb.NewMessage(update.Message.Chat.ID, "El servidor no responde :(")
		bot.Send(msg)

		return err
	}

	var sb strings.Builder

	// blah blah blah
	for i := range news {
		if i >= limit {
			break
		}

		sb.WriteString("> ")
		sb.WriteString(news[i].Anchor)
		sb.WriteRune('\n')
	}

	msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
	msg.ParseMode = "html"
	_, err = bot.Send(msg)

	return err
}

func NewsCmd(ctx commander.Context) error {
	return handleCmd(ctx, "news", 5)
}

func AvisosCmd(ctx commander.Context) error {
	return handleCmd(ctx, "avisos", 5)
}

func CoreCmd(ctx commander.Context) error {
	return handleCmd(ctx, "coredumped", 3)
}
