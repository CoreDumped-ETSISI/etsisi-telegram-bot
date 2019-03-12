package help

import (
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func HelpCmd(ctx commander.Context) error {
	bot := ctx.Arg("bot").(*tb.BotAPI)
	update := ctx.Arg("update").(tb.Update)

	text := `Comandos disponibles:

/menu - Menú de la cafetería de hoy.
/menu2 - Menú de la cafetería de mañana.
/salas - Disponibilidad de salas de trabajo de la biblioteca.
/noticias - Noticias generales de ETSISI.
/avisos - Avisos para los alumnos.
/coredumped - Blog de la asociación CoreDumped.
/subscribe (canal?) - Suscribirse a un canal (noticias, avisos o status). Si no se especifica, se suscribe a todos los canales.
/unsubscribe (canal?) - Desuscribirse de un canal (noticias, avisos o status). Si no se especifica, se desuscribe de todos los canales.
/horario (grupo?) - Horario para hoy de tu grupo.
/horario2 (grupo?) - Horario de la semana de tu grupo.
/status - Comprueba el estado de varios servicios de la uni.
/help - Este mensaje.
`

	msg := tb.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "markdown"

	_, err := bot.Send(msg)

	return err
}
