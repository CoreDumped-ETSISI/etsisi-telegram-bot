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
	callbacks := commander.New()

	me, err := bot.GetMe()

	if err != nil {
		panic(err)
	}

	cmd.Preprocessor = &CustomTelegramPreprocessor{
		BotName: me.UserName,
	}

	route(cmd, config, callbacks)
	use(cmd, config)

	for update := range updates {
		go handleUpdate(config, bot, update, cmd, callbacks)
	}
}

func handleUpdate(state state.T, bot *tb.BotAPI, update tb.Update, cmd *commander.CommandGroup, callbacks *commander.CommandGroup) {
	if update.Message != nil && update.Message.Text != "" {
		ok, err := cmd.ExecuteWithContext(update.Message.Text, map[string]interface{}{
			"bot":    bot,
			"update": update,
			"state": state,
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

		cmd.TriggerWithContext("text", map[string]interface{}{
			"bot":    bot,
			"update": update,
			"state": state,
		})
	} else if update.CallbackQuery != nil {
		cq := update.CallbackQuery

		ok, err := callbacks.ExecuteWithContext(cq.Data, map[string]interface{}{
			"bot":    bot,
			"update": update,
			"state": state,
		})

		log.WithFields(log.Fields{
			"chatid": cq.Message.Chat.ID,
			"chat":   getChatTitle(cq.Message),
			"sender": getSenderName(cq.Message.From),
			"data":   cq.Data,
			"error":  err,
			"found":  ok,
		}).Debug("Got query callback")
	}

}
