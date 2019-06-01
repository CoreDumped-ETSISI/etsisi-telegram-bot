package guides

import (
	"fmt"
	"net/http"
	"strings"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func GuideCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	g, err := getAllGuides()

	if err != nil {
		return err
	}

	msg := tb.NewMessage(update.Message.Chat.ID, "Por favor, seleccione el grado.")

	var blist [][]tb.InlineKeyboardButton

	for grado := range g {
		button := tb.NewInlineKeyboardButtonData(grado, fmt.Sprintf("/gpag %v 0", grado))
		row := tb.NewInlineKeyboardRow(button)
		blist = append(blist, row)
	}

	markup := tb.NewInlineKeyboardMarkup(blist...)

	msg.ReplyMarkup = markup

	_, err = bot.Send(msg)

	return err
}

func DownloadGuideCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)
	code := ctx.ArgString("code")

	g, err := getAllGuides()

	if err != nil {
		return err
	}

	var guia *Guide

outlp:
	for grado := range g {
		for i := range g[grado] {
			if g[grado][i].Code == code {
				guia = g[grado][i]
				break outlp
			}
		}
	}

	if guia == nil {
		return nil
	}

	r, err := http.Get(guia.URL)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	doc := tb.FileReader{
		Name:   guia.Name + ".pdf",
		Reader: r.Body,
		Size:   -1,
	}

	msg := tb.NewDocumentUpload(update.Message.Chat.ID, doc)
	msg.Caption = fmt.Sprintf("📘<b>%v</b>\nSemestre: %v\nTipo: %v\nCréditos: %v ECTS", guia.Name, guia.Semester, guia.Type, guia.ECTS)
	msg.ParseMode = "html"

	_, err = bot.Send(msg)

	return err
}

// Callback /gpag {grado} {offset}
func PaginateGradoCallback(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)
	grado := ctx.ArgString("grado")
	offset := ctx.ArgInt("offset")

	originalMsg := update.CallbackQuery.Message

	g, err := getAllGuides()

	if err != nil {
		return err
	}

	var sb strings.Builder

	fin := offset + 5
	l := len(g[grado])

	if fin > l {
		fin = l
	}

	for _, guia := range g[grado][offset:fin] {
		sb.WriteString(writeGuide(guia))
	}

	edit := tb.NewEditMessageText(originalMsg.Chat.ID, originalMsg.MessageID, sb.String())
	edit.ParseMode = "html"

	// Pagination

	var blist []tb.InlineKeyboardButton

	if offset > 0 {
		left := tb.NewInlineKeyboardButtonData("<<", fmt.Sprintf("/gpag %v %v", grado, offset-5))
		blist = append(blist, left)
	}

	if fin < len(g[grado]) {
		right := tb.NewInlineKeyboardButtonData(">>", fmt.Sprintf("/gpag %v %v", grado, offset+5))
		blist = append(blist, right)
	}

	row := tb.NewInlineKeyboardRow(blist...)
	markup := tb.NewInlineKeyboardMarkup(row)

	edit.ReplyMarkup = &markup

	_, _ = bot.AnswerCallbackQuery(tb.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
		Text:            "",
	})

	_, err = bot.Send(edit)

	return err
}

func writeGuide(g *Guide) string {
	return fmt.Sprintf("📘 <b>%v</b>\nSemestre %v | %v\nDescargar: /gg_%v\n\n", g.Name, g.Semester, g.Type, g.Code)
}
