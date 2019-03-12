package horario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func HorarioCmd(redis *redis.Client) func(commander.Context) error {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)

		grupo := ctx.ArgString("grupo")
		oldfav := false

		if grupo == "" {
			grp, err := redis.Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

			if err != nil {
				// Has no favorite group and didn't provide a group.
				msg := tb.NewMessage(update.Message.Chat.ID, "No tienes un grupo favorito. Especifica tu grupo con /horario (GRUPO)")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			grupo = grp
			oldfav = true
		}

		horario, err := getHorarioForGroup(grupo)

		if err != nil {
			if err == NoSuchGroupError {
				msg := tb.NewMessage(update.Message.Chat.ID, "Ese grupo no existe")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			return err
		}

		if !oldfav {
			err = redis.Set(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID), grupo, 0).Err()

			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"grupo":  grupo,
					"userid": update.Message.From.ID,
				}).Error("Error when setting favorite group")
			}
		}

		now := int(time.Now().UTC().Weekday()) - 1
		if now < 0 {
			now += 7
		}

		if now >= len(horario) || len(horario[now]) == 0 {
			msg := tb.NewMessage(update.Message.Chat.ID, "Hoy no tienes clase")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			return nil
		}

		horarioToday := horario[now]

		var sb strings.Builder

		sb.WriteString("Tu horario para hoy:\n\n")

		for i := range horarioToday {
			clase := horarioToday[i]

			sb.WriteString(fmt.Sprintf("%v:00-%v:00 tienes %v\n", clase.Start, clase.End, clase.Name))
		}

		msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
		_, err = bot.Send(msg)

		return err
	}
}

func HorarioWeekCmd(redis *redis.Client) func(commander.Context) error {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)

		grupo := ctx.ArgString("grupo")
		oldfav := false

		if grupo == "" {
			grp, err := redis.Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

			if err != nil {
				// Has no favorite group and didn't provide a group.
				msg := tb.NewMessage(update.Message.Chat.ID, "No tienes un grupo favorito. Especifica tu grupo con /horario2 (GRUPO)")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			grupo = grp
			oldfav = true
		}

		horario, err := getHorarioForGroup(grupo)

		if err != nil {
			if err == NoSuchGroupError {
				msg := tb.NewMessage(update.Message.Chat.ID, "Ese grupo no existe")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				return nil
			}

			return err
		}

		if !oldfav {
			err = redis.Set(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID), grupo, 0).Err()

			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"grupo":  grupo,
					"userid": update.Message.From.ID,
				}).Error("Error when setting favorite group")
			}
		}

		// ASCII tables don't fit in group conversations.
		if !update.Message.Chat.IsPrivate() {
			var buf bytes.Buffer
			err = json.NewEncoder(&buf).Encode(horario)

			img, err := http.Post("https://renderer.kolhos.chichasov.es/api/horario", "application/json", &buf)

			if err != nil {
				return err
			}

			defer img.Body.Close()

			file := tb.FileReader{
				Name:   "horario.png",
				Size:   -1,
				Reader: img.Body,
			}

			msg := tb.NewPhotoUpload(update.Message.Chat.ID, file)
			msg.Caption = "Horario de " + strings.ToUpper(grupo)
			_, err = bot.Send(msg)
			return err
		} else {
			tableHeader := []string{"⌚", "M", "X", "J", "V"}
			tableData := [][]string{}
			hasRow := []bool{}

			// Generate table
			for i := 9; i <= 20; i++ {
				tableData = append(tableData, []string{fmt.Sprintf("%v", i), "", "", "", ""})
				hasRow = append(hasRow, false)
			}

			// Populate table
			// Start without monday or it won't fit on mobile
			for day := 1; day < len(horario); day++ {
				for _, clase := range horario[day] {
					for i := clase.Start; i < clase.End; i++ {
						tableData[i-9][day] = clase.Name
						hasRow[i-9] = true
					}
				}
			}

			// Expunge all rows without data
			for i := len(hasRow) - 1; i >= 0; i-- {
				if !hasRow[i] {
					copy(tableData[i:], tableData[i+1:])
					tableData[len(tableData)-1] = nil
					tableData = tableData[:len(tableData)-1]
				}
			}

			var sb strings.Builder
			sb.WriteString("<pre>\n")
			table := tablewriter.NewWriter(&sb)

			table.SetHeader(tableHeader)
			table.SetAutoMergeCells(true)
			table.SetRowLine(true)
			table.AppendBulk(tableData)
			table.SetCaption(true, fmt.Sprintf("Horario de %v", strings.ToUpper(grupo)))

			table.Render()

			sb.WriteString("</pre>")

			msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
			msg.ParseMode = "html"
			_, err = bot.Send(msg)

			return err
		}
	}
}
