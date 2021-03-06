package subscription

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func SubscribeCmd(s *DBContext) func(commander.Context) error {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)
		bot := update.State.Bot()
		feed := ctx.ArgString("feed")

		chatid := update.Message.Chat.ID
		if channel, ok := normalMap[feed]; ok {
			key := channel + "_SUBSCRIBERS"
			err := s.addSubscriber(key, chatid)

			if err != nil {
				msg := tb.NewMessage(chatid, "No se ha podido suscribirse :(")
				bot.Send(msg)
				return err
			}

			msg := tb.NewMessage(chatid, "Se ha suscrito al canal `"+feed+"`")
			msg.ParseMode = "markdown"
			bot.Send(msg)
		} else if feed == "" {
			for _, key := range publicChannels {
				err := s.addSubscriber(normalMap[key]+"_SUBSCRIBERS", chatid)

				if err != nil {
					msg := tb.NewMessage(chatid, "No se ha podido suscribirse :(")
					bot.Send(msg)
					return err
				}
			}

			msg := tb.NewMessage(chatid, "Se ha suscrito a todos los canales!")
			bot.Send(msg)
		} else {
			msg := tb.NewMessage(chatid, "Ese canal no existe!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		return nil
	}
}

func UnsubscribeCmd(s *DBContext) func(commander.Context) error {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)
		bot := update.State.Bot()
		feed := ctx.ArgString("feed")

		chatid := update.Message.Chat.ID
		if channel, ok := normalMap[feed]; ok {
			err := s.removeSubscriber(channel+"_SUBSCRIBERS", chatid)

			if err != nil {
				msg := tb.NewMessage(chatid, "No se ha podido desuscribirse :(")
				bot.Send(msg)
				return err
			}

			msg := tb.NewMessage(chatid, "Se ha desuscrito del canal `"+feed+"`")
			msg.ParseMode = "markdown"
			bot.Send(msg)
		} else if feed == "" {
			for key := range normalMap {
				err := s.removeSubscriber(normalMap[key]+"_SUBSCRIBERS", chatid)

				if err != nil {
					msg := tb.NewMessage(chatid, "No se ha podido desuscribirse :(")
					bot.Send(msg)
					return err
				}
			}

			msg := tb.NewMessage(chatid, "Se ha desuscrito de todos los canales!")
			bot.Send(msg)
		} else {
			msg := tb.NewMessage(chatid, "Ese canal no existe!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		return nil
	}
}

func StartMonitoringSubscriptions(redis *redis.Client, bot *tb.BotAPI, s *DBContext) {
	pubsub := redis.Subscribe(redisChannels...)

	ch := pubsub.Channel()

	limiter := time.Tick(10 * time.Minute)

	for msg := range ch {
		select {
		case <-limiter:
			// If we're under the rate-limit, proceed
			break
		default:
			// Otherwise drop the event
			continue
		}

		var item channelMessage
		err := json.Unmarshal([]byte(msg.Payload), &item)

		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"channel": msg.Channel,
					"payload": msg.Payload,
				}).Error("Error deserializing redis message")
		} else {
			log.
				WithFields(log.Fields{
					"channel": msg.Channel,
					"payload": msg.Payload,
				}).Info("Received event on channel")

			members, err := s.getSubscribers(msg.Channel + "_SUBSCRIBERS")

			if err != nil {
				log.
					WithError(err).
					WithFields(log.Fields{
						"channel": msg.Channel,
						"payload": msg.Payload,
					}).Error("Error listing suscription members")
			} else {
				humanChannel := reverseMap[msg.Channel]

				var sb strings.Builder
				sb.WriteString("<b>Nuevo mensaje en el canal </b><i>")
				sb.WriteString(humanChannel)
				sb.WriteString("</i>\n")
				sb.WriteString(item.Text)

				var markup tb.InlineKeyboardMarkup

				if item.Link != nil {
					button := tb.NewInlineKeyboardButtonURL("Más Info", *item.Link)
					markup = tb.NewInlineKeyboardMarkup([]tb.InlineKeyboardButton{button})
				}

				text := sb.String()

				for i := range members {
					memberid := members[i]

					message := tb.NewMessage(memberid, text)
					message.ParseMode = "html"
					if item.Link != nil {
						message.ReplyMarkup = markup
					}
					_, err = bot.Send(message)

					if err != nil {
						log.
							WithError(err).
							WithFields(log.Fields{
								"channel": msg.Channel,
								"payload": msg.Payload,
								"userid":  memberid,
							}).Error("Error sending news item to user")
					}
				}

				log.
					WithFields(log.Fields{
						"channel":     msg.Channel,
						"subscribers": len(members),
					}).Info("Sent to subscribers")
			}
		}
	}
}

func GetAllChannelsCommand(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)
	bot := update.State.Bot()

	var sb strings.Builder

	sb.WriteString("Canales disponibles:\n\n")

	for key := range normalMap {
		sb.WriteString(key)
		sb.WriteRune('\n')
	}

	msg := tb.NewMessage(update.Message.Chat.ID, sb.String())
	_, err := bot.Send(msg)
	return err
}
