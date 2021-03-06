package verify

import (
	"errors"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

var (
	ErrNotPrivateChat = errors.New("This must be a private chat")
)

func PrivateOnlyMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		var msg *tb.Message

		// Text messages
		if update.Message != nil {
			msg = update.Message
		} else if update.CallbackQuery != nil {
			// Callbacks
			msg = update.CallbackQuery.Message
		} else {
			return ErrNotPrivateChat
		}

		if !msg.Chat.IsPrivate() {
			return ErrNotPrivateChat
		}

		return next(ctx)
	}
}
