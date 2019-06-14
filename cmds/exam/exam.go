package exam

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func getUserGrupo(redis *redis.Client, update tb.Update) string {
	grp, err := redis.Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

	if err != nil {
		return ""
	}

	return grp
}

func getTagsForGroup(grp string) [][]string {
	grp = strings.ToUpper(grp)

	shared1 := regexp.MustCompile(`G[MT]1\d`)
	shared2 := regexp.MustCompile(`G[MT]2\d`)

	if shared1.MatchString(grp) {
		return [][]string{[]string{"general", "year_1"}}
	}

	if shared2.MatchString(grp) {
		return [][]string{[]string{"general", "year_2"}}
	}

	comp := regexp.MustCompile(`GCO[MT]3\d`)
	soft := regexp.MustCompile(`GIW[TM]3\d`)
	si := regexp.MustCompile(`GSI[TM]3\d`)
	ti := regexp.MustCompile(`GTI[TM]3\d`)

	if comp.MatchString(grp) {
		return [][]string{[]string{"general", "year_3"},
			[]string{"compu", "year_3"},
		}
	}

	if soft.MatchString(grp) {
		return [][]string{[]string{"general", "year_3"},
			[]string{"software", "year_3"},
		}
	}

	if si.MatchString(grp) {
		return [][]string{[]string{"general", "year_3"},
			[]string{"si", "year_3"},
		}
	}

	if ti.MatchString(grp) {
		return [][]string{[]string{"general", "year_3"},
			[]string{"ti", "year_3"},
		}
	}

	last := regexp.MustCompile(`4\d$`)

	if last.MatchString(grp) {
		return [][]string{
			[]string{"year_4"},
		}
	}

	return nil
}

func ExamCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)
	state := ctx.Arg("state").(state.T)

	ug := getUserGrupo(state.Redis(), update)

	if ug != "" {
		tags := getTagsForGroup(ug)

		if tags != nil {
			ctx.AddArg("tags", tags)
			return ShowExamsCb(ctx)
		}
	}

	msg := tb.NewMessage(update.Message.Chat.ID, "Por favor, seleccione el grado.")

	var blist [][]tb.InlineKeyboardButton

	for k, v := range Grados {
		button := tb.NewInlineKeyboardButtonData(k, fmt.Sprintf("/exyear %v", v))
		row := tb.NewInlineKeyboardRow(button)
		blist = append(blist, row)
	}

	markup := tb.NewInlineKeyboardMarkup(blist...)

	msg.ReplyMarkup = markup

	_, err := bot.Send(msg)

	return err
}

// /exyear {grado}
func SelectYearCb(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	grado := ctx.ArgString("grado")

	msg := tb.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Por favor, seleccione el curso.")

	var blist [][]tb.InlineKeyboardButton

	for i := 1; i <= 4; i++ {
		button := tb.NewInlineKeyboardButtonData(fmt.Sprintf("Curso %v", i), fmt.Sprintf("/exshow %v %v", grado, i))
		row := tb.NewInlineKeyboardRow(button)
		blist = append(blist, row)
	}

	markup := tb.NewInlineKeyboardMarkup(blist...)

	msg.ReplyMarkup = &markup

	_, err := bot.Send(msg)

	return err
}

// /exshow {grado} {curso:int}
func ShowExamsCb(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	grado := ctx.ArgString("grado")
	curso := ctx.ArgInt("curso")
	cursotag := fmt.Sprintf("year_%v", curso)

	ex, err := getAllExams()

	if err != nil {
		return err
	}

	params, ok := ctx.Arg("tags").([][]string)

	if ok && params != nil {
		log.WithField("tags", params).Debug("found tags for group")

		ex = filterByTags(ex, params...)
	} else {
		if curso <= 2 {
			ex = filterByTags(ex, []string{cursotag, "general"})
		} else {
			ex = filterByTags(ex, []string{cursotag, grado}, []string{cursotag, "general"})
		}
	}

	// TODO: Dont hardcode this
	extraordinaria := time.Date(2019, time.June, 17, 0, 0, 0, 0, time.Local)

	if time.Now().After(extraordinaria) {
		ex = filterByDate(ex, time.Now(), time.Now().AddDate(1, 0, 0))
	} else {
		ex = filterByDate(ex, time.Now(), extraordinaria)
	}

	var sb strings.Builder

	for _, exam := range ex {
		sb.WriteString(fmt.Sprintf("📚 %v - <b>%v</b> (%v en Bloque %v)\n", exam.Day, exam.Name, exam.Timeslot, strings.Join(exam.Aulas, "/")))
	}

	if sb.Len() == 0 {
		sb.WriteString("Por ahora no tienes examenes.")
	}

	if params != nil {
		msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
		msg.ParseMode = "html"
		_, err = bot.Send(msg)
	} else {
		msg := tb.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, sb.String())
		msg.ParseMode = "html"
		_, err = bot.Send(msg)
	}

	return err
}
