package main

import (
	"math/rand"
	"net/http"
	"os"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func startBot() (*tb.BotAPI, tb.UpdatesChannel) {
	bot, err := tb.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))

	if err != nil {
		panic(err)
	}

	url := os.Getenv("WEBHOOK_URL")

	if url == "" {
		// No webhook setup, use polling
		_, _ = bot.RemoveWebhook()
		u := tb.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)

		if err != nil {
			panic(err)
		}

		return bot, updates
	}

	// Generate random webhook URL
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	n := 10
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	token := string(b)

	_, err = bot.SetWebhook(tb.NewWebhook("https://" + url + "/wh/" + token))

	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/wh/" + token)

	port := "8080"

	if ep := os.Getenv("WEBHOOK_PORT"); ep != "" {
		port = ep
	}

	go http.ListenAndServe("0.0.0.0:"+port, nil)

	log.Info("Iniciado y escuchando!")

	return bot, updates
}
