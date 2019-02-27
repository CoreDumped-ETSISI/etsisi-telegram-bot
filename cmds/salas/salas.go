package salas

import (
	"io"
	"time"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func SalasCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	chatID := update.Message.Chat.ID

	bib, err := getSalas()

	if err != nil {
		return err
	}

	now := time.Now().In(time.UTC).Add(1 * time.Hour)

	r, w := io.Pipe()

	go generateImage(bib, timeToIndex(now.Hour(), now.Minute())+1, w)

	file := tb.FileReader{
		Name:   "Salas de trabajo.png",
		Size:   -1,
		Reader: r,
	}

	msg := tb.NewPhotoUpload(chatID, file)

	_, err = bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
