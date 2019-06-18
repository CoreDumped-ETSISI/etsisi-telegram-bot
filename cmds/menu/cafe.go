package menu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func CafeTodayCmd(ctx commander.Context) error {
	return sendMenuToChat(ctx, "today")
}

func CafeTomorrowCmd(ctx commander.Context) error {
	return sendMenuToChat(ctx, "tomorrow")
}

func sendMenuToChat(ctx commander.Context, endpoint string) error {
	update := ctx.Arg("update").(state.Update)
	bot := update.State.Bot()

	chatID := update.Message.Chat.ID

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://cafe.kolhos.chichasov.es/%s", endpoint))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 { // Mensaje de error.
		var errorMsg ErrorMessage
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)

		if err != nil {
			return err
		}

		msg := tb.NewMessage(chatID, fmt.Sprintf("Ups, algo ha ido mal:\n%v", errorMsg.Message))

		bot.Send(msg)

		return nil
	}

	var menu MenuDia
	err = json.NewDecoder(resp.Body).Decode(&menu)

	if err != nil {
		return err
	}

	sb := strings.Builder{}

	dia := "Hoy"

	if endpoint == "tomorrow" {
		dia = "Mañana"
	}

	sb.WriteString(dia)
	sb.WriteString(" de primer plato hay:\n")

	for _, plato := range menu.PrimerPlato {
		sb.WriteString("- ")
		sb.WriteString(plato)
		sb.WriteRune('\n')
	}

	sb.WriteRune('\n')
	sb.WriteString(dia)
	sb.WriteString(" de segundo plato hay:\n")

	for _, plato := range menu.SegundoPlato {
		sb.WriteString("\t- ")
		sb.WriteString(plato)
		sb.WriteRune('\n')
	}

	msg := tb.NewMessage(chatID, sb.String())

	_, err = bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
