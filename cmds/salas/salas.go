package salas

import (
	"net/http"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func SalasCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	chatID := update.Message.Chat.ID

	resp, err := http.Get("https://biblio.kolhos.chichasov.es/api/salas")

	if err != nil {
		return err
	}

	img, err := http.Post("https://renderer.kolhos.chichasov.es/api/salas", "application/json", resp.Body)

	if err != nil {
		return err
	}

	defer img.Body.Close()

	file := tb.FileReader{
		Name:   "Salas de trabajo.png",
		Size:   -1,
		Reader: img.Body,
	}

	msg := tb.NewPhotoUpload(chatID, file)

	_, err = bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
