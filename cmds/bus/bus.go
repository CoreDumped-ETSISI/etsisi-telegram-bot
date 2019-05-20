package bus

import (
	"fmt"
	"strings"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func BusCmd(ctx commander.Context) error {
	stopid := ctx.ArgInt("stop")

	if stopid == 0 {
		return uniBusCmd(ctx)
	} else {
		return busStopCmd(ctx)
	}
}

func busStopCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	stopid := ctx.ArgInt("stop")

	arrives, err := getEstimatesForStop(stopid)

	var sb strings.Builder

	sb.WriteString("<b>Próximas Llegadas</b>\n")

	// 🚌 E - CONDE DE CASAL (2m)
	for _, bus := range arrives {
		sb.WriteString("🚌 ")
		sb.WriteString(bus.LineID)
		sb.WriteString(" - ")
		sb.WriteString(bus.Destination)
		sb.WriteString(fmt.Sprintf(" (%vm)\n", int(bus.TimeLeft.Minutes())))
	}

	msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
	msg.ParseMode = "html"
	_, err = bot.Send(msg)
	return err
}

func uniBusCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	loadingMsg := tb.NewMessage(update.Message.Chat.ID, "Cargando ⏰")
	loadingMsg.ParseMode = "html"
	lmsg, err := bot.Send(loadingMsg)

	if err != nil {
		return err
	}

	arrives, err := getUniEstimates()

	if err != nil {
		m := tb.NewEditMessageText(update.Message.Chat.ID, lmsg.MessageID, "Algo salió mal 💀")
		_, _ = bot.Send(m)

		return err
	}

	var sb strings.Builder

	sb.WriteString("<b>Paradas ETSISI</b>\n")
	sb.WriteString("<b>Sentido Conde de Casal</b> (#4281)\n")

	// 🚌 E - 2m0s -
	for _, bus := range arrives.SentidoConde {
		sb.WriteString("🚌 ")
		sb.WriteString(bus.LineID)
		sb.WriteString(fmt.Sprintf(" - %v - %vm\n", bus.TimeLeft, bus.Distance))
	}

	sb.WriteString("<b>Sentido Sierra</b> (#4702)\n")

	// 🚌 E - CONDE DE CASAL (2m)
	for _, bus := range arrives.SentidoSierra {
		sb.WriteString("🚌 ")
		sb.WriteString(bus.LineID)
		sb.WriteString(fmt.Sprintf(" - %v - %vm\n", bus.TimeLeft, bus.Distance))
	}

	sb.WriteString("<b>Puente Gasolinera</b> (#2712)\n")

	for _, bus := range arrives.Puente {
		sb.WriteString("🚌 ")
		sb.WriteString(bus.LineID)
		sb.WriteString(fmt.Sprintf(" - %v - %vm\n", bus.TimeLeft, bus.Distance))
	}

	msg := tb.NewEditMessageText(update.Message.Chat.ID, lmsg.MessageID, sb.String())
	msg.ParseMode = "html"
	_, err = bot.Send(msg)
	return err
}
