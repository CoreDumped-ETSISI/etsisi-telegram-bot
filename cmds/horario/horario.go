package horario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

// A lot of people misinterpret instructions
// and do e.g. /horario (GM11) instead of
// /horario GM11
func cleanGroupArg(grp string) string {
	return strings.Trim(grp, "()")
}

func HorarioCmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)

	grupo := ctx.ArgString("grupo")
	grupo = cleanGroupArg(grupo)
	oldfav := false

	if grupo == "" {
		grp, err := update.State.Redis().Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

		if err != nil {
			// Has no favorite group and didn't provide a group.
			msg := tb.NewMessage(update.Message.Chat.ID, "No tienes un grupo favorito. Especifica tu grupo con /horario (GRUPO)")
			msg.ReplyToMessageID = update.Message.MessageID
			update.State.Bot().Send(msg)
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
			update.State.Bot().Send(msg)
			return nil
		}

		return err
	}

	if !oldfav {
		err = update.State.Redis().Set(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID), grupo, 0).Err()

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
		update.State.Bot().Send(msg)
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
	_, err = update.State.Bot().Send(msg)

	return err
}

func HorarioWeekCmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)
	bot := update.State.Bot()

	grupo := ctx.ArgString("grupo")
	grupo = cleanGroupArg(grupo)
	oldfav := false

	if grupo == "" {
		grp, err := update.State.Redis().Get(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID)).Result()

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
		err = update.State.Redis().Set(fmt.Sprintf("FAVORITE_GROUP_%v", update.Message.From.ID), grupo, 0).Err()

		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"grupo":  grupo,
				"userid": update.Message.From.ID,
			}).Error("Error when setting favorite group")
		}
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(horario)
	
	if err != nil {
		return err
	}

	img, err := http.Post(services.Get("renderer",8080)+"/api/horario", "application/json", &buf)

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
}
