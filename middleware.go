package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
	log "github.com/sirupsen/logrus"
)

func (cfg config) ratelimitMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		key := fmt.Sprintf("TIMEOUT_%v_%v", update.Message.Chat.ID, ctx.Name)
		_, err := cfg.redis.Get(key).Result()

		if err == nil { // Key was found
			sendNag(update.State.Bot(), update.Message.Chat.ID)
			return nil
		}

		err = cfg.redis.Set(key, true, cfg.commandTimeout).Err()

		if err != nil {
			// If redis is unavailable, don't execute the command to prevent spam.
			return err
		}

		return next(ctx)
	}
}

func sendNag(bot *tb.BotAPI, chatID int64) {
	frases := []string{
		"Deja de hacer spam, prostifruto!",
		"El spam a tu casa, campeón.",
		"Stop spamming, idiota. ¿Tanta prisa tienes?",
	}

	msg := tb.NewMessage(chatID, frases[rand.Intn(len(frases))])

	bot.Send(msg)
}

func loggerMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		pre := time.Now()

		// TODO: Recover from panic
		err := next(ctx)

		elapsed := time.Now().Sub(pre)

		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"command": ctx.Name,
					"args":    ctx.Args,
					"chatid":  update.Message.Chat.ID,
					"chat":    getChatTitle(update.Message),
					"sender":  getSenderName(update.Message.From),
					"text":    update.Message.Text,
					"elapsed": elapsed,
				}).Error("Error when executing command")
		} else {
			log.WithFields(log.Fields{
				"command": ctx.Name,
				"args":    ctx.Args,
				"chatid":  update.Message.Chat.ID,
				"chat":    getChatTitle(update.Message),
				"sender":  getSenderName(update.Message.From),
				"text":    update.Message.Text,
				"elapsed": elapsed,
			}).Info("Successfuly sent command")
		}

		return err
	}
}

func callbackLoggerMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		pre := time.Now()

		// TODO: Recover from panic
		err := next(ctx)

		elapsed := time.Now().Sub(pre)

		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"command": ctx.Name,
					"args":    ctx.Args,
					"chatid":  update.CallbackQuery.Message.Chat.ID,
					"chat":    getChatTitle(update.CallbackQuery.Message),
					"sender":  getSenderName(update.CallbackQuery.Message.From),
					"text":    update.CallbackQuery.Data,
					"elapsed": elapsed,
				}).Error("Error when executing command")
		}

		return err
	}
}

func use(cmd *commander.CommandGroup, cfg config) {
	cmd.Use(loggerMiddleware)
	cmd.Use(cfg.ratelimitMiddleware)
}
