package exam

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func getUserGrupo(redis *redis.Client, update tb.Update) string {
	grp, err := redis.Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

	if err != nil {
		return ""
	}

	return grp
}

func getTagsForGroup(grp string) [][]string {
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
		return [][]string{[]string{"general", "year_3"},
			[]string{"year_4"},
		}
	}

	return nil
}

func ExamCmd(redis *redis.Client) func(commander.Context) error {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)
		params := ctx.ArgString("params")

		ex, err := getAllExams()

		if err != nil {
			return err
		}

		// TODO: Dont hardcode this
		extraordinaria := time.Date(2019, time.June, 17, 0, 0, 0, 0, time.Local)

		if time.Now().After(extraordinaria) {
			ex = filterByDate(ex, time.Now(), time.Now().AddDate(1, 0, 0))
		} else {
			ex = filterByDate(ex, time.Now(), extraordinaria)
		}

		if params != "" {
			tags := strings.Split(params, " ")
			ex = filterByTags(ex, tags)
		} else {
			ug := getUserGrupo(redis, update)

			if ug == "" {
				msg := tb.NewMessage(update.Message.Chat.ID, getHelpMsg())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			tags := getTagsForGroup(ug)

			if tags == nil {
				msg := tb.NewMessage(update.Message.Chat.ID, getHelpMsg())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			ex = filterByTags(ex, tags...)
		}

		var sb strings.Builder

		for _, exam := range ex {
			sb.WriteString(fmt.Sprintf("📚 %v - <b>%v</b> (%v en Bloque %v)\n", exam.Day, exam.Name, exam.Timeslot, strings.Join(exam.Aulas, "/")))
		}

		msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "html"
		bot.Send(msg)

		return nil
	}
}

func getHelpMsg() string {
	return "No sé de que grupo eres. Puedes guardar tu grupo con /horario (tu grupo) o buscar examenes con tags, e.g.\n/exam software year_1 optativa\n/exam ti year_4"
}
