package help

import (
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/verify"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func Start(ctx commander.Context) error {
	data := ctx.Arg("data")

	if data == "" {
		return HelpCmd(ctx)
	}

	switch data {
	case "verifyme":
		return verify.Cmd(ctx)
	}

	return nil
}

func HelpCmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)

	text := `Comandos disponibles:

/menu - Menú de la cafetería de hoy.
/menu2 - Menú de la cafetería de mañana.
/salas - Disponibilidad de salas de trabajo de la biblioteca.
/noticias - Noticias generales de ETSISI.
/avisos - Avisos para los alumnos.
/coredumped - Blog de la asociación CoreDumped.
/subscribe (canal?) - Suscribirse a un canal (/canales). Si no se especifica, se suscribe a todos los canales.
/unsubscribe (canal?) - Desuscribirse de un canal (/canales). Si no se especifica, se desuscribe de todos los canales.
/horario (grupo?) - Horario para hoy de tu grupo.
/horario2 (grupo?) - Horario de la semana de tu grupo.
/exam  - Exámenes finales de tu curso.
/guias - Descargar guías de asignaturas.
/status - Comprueba el estado de varios servicios de la uni.
/bus (id parada?) - Tiempo de buses de la uni o de una parada concreta.
/help - Este mensaje.
`

	msg := tb.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "markdown"

	_, err := update.State.Bot().Send(msg)

	return err
}
