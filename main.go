package main

import (
	"github.com/guad/commander"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := newConfig()

	log.SetLevel(log.AllLevels[config.logLevel])

	bot, updates := startBot()

	config.bot = bot

	cmd := commander.New()

	me, err := bot.GetMe()

	if err != nil {
		panic(err)
	}

	cmd.Preprocessor = &commander.TelegramPreprocessor{
		BotName: me.UserName,
	}

	route(cmd, config)
	use(cmd, config)

	for update := range updates {
		go handleUpdate(bot, update, cmd)
	}
}

func handleUpdate(bot *tb.BotAPI, update tb.Update, cmd *commander.CommandGroup) {
	if update.Message != nil && update.Message.Text != "" {
		ok, err := cmd.ExecuteWithContext(update.Message.Text, map[string]interface{}{
			"bot":    bot,
			"update": update,
		})

		// TODO: Maybe send the user a message with the error?

		log.WithFields(log.Fields{
			"chatid": update.Message.Chat.ID,
			"chat":   getChatTitle(update.Message),
			"sender": getSenderName(update.Message.From),
			"text":   update.Message.Text,
			"error":  err,
			"found":  ok,
		}).Debug("Got update")

		// General logging is done by the logging middleware.
	}
}
