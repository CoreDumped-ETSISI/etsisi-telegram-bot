package subscription

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func SubscribeCmd(redis *redis.Client) func(commander.Context) error {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)
		feed := ctx.ArgString("feed")

		chatid := update.Message.Chat.ID
		if channel, ok := normalMap[feed]; ok {
			_, err := redis.SAdd(channel+"_SUBSCRIBERS", chatid).Result()

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
				_, err := redis.SAdd(normalMap[key]+"_SUBSCRIBERS", chatid).Result()

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

func UnsubscribeCmd(redis *redis.Client) func(commander.Context) error {
	return func(ctx commander.Context) error {
		bot := ctx.Arg("bot").(*tb.BotAPI)
		update := ctx.Arg("update").(tb.Update)
		feed := ctx.ArgString("feed")

		chatid := update.Message.Chat.ID
		if channel, ok := normalMap[feed]; ok {
			_, err := redis.SRem(channel+"_SUBSCRIBERS", chatid).Result()

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
				_, err := redis.SRem(normalMap[key]+"_SUBSCRIBERS", chatid).Result()

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

func StartMonitoringSubscriptions(redis *redis.Client, bot *tb.BotAPI) {
	pubsub := redis.Subscribe(redisChannels...)

	ch := pubsub.Channel()

	for msg := range ch {
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

			members, err := redis.SMembers(msg.Channel + "_SUBSCRIBERS").Result()

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
					memberid, _ := strconv.ParseInt(members[i], 10, 64)

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
			}
		}
	}
}
