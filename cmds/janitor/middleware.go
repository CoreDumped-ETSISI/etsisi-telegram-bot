package janitor

import (
	"errors"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

var (
	ErrNotAnAdmin     = errors.New("You are not an administrator")
	ErrChatNotManaged = errors.New("This chat is not managed")
)

func AdminOnlyMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		var chatid int64
		var userid int

		// Text messages
		if update.Message != nil {
			// Do nothing in private chats
			if update.Message.Chat.IsPrivate() {
				return next(ctx)
			}
			chatid = update.Message.Chat.ID
			userid = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			// Callbacks
			// Do nothing in private chats
			if update.CallbackQuery.Message.Chat.IsPrivate() {
				return next(ctx)
			}

			chatid = update.CallbackQuery.Message.Chat.ID
			userid = update.CallbackQuery.From.ID
		} else {
			return nil
		}

		m, err := update.State.Bot().GetChatMember(tb.ChatConfigWithUser{
			ChatID: chatid,
			UserID: userid,
		})

		// Deny if something went wrong.
		if err != nil {
			return err
		}

		if m.IsAdministrator() || m.IsCreator() {
			return next(ctx)
		}

		return ErrNotAnAdmin
	}
}

func ManagedOnlyMiddleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		update := ctx.Arg("update").(state.Update)

		var chatid int64

		// Text messages
		if update.Message != nil {
			// Do nothing in private chats
			if update.Message.Chat.IsPrivate() {
				return next(ctx)
			}
			chatid = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			// Callbacks
			// Do nothing in private chats
			if update.CallbackQuery.Message.Chat.IsPrivate() {
				return next(ctx)
			}

			chatid = update.CallbackQuery.Message.Chat.ID
		} else {
			return nil
		}

		managed, err := isChatManaged(chatid)

		if err != nil {
			return err
		}

		// If it's not managed, stop.
		if !managed {
			return ErrChatNotManaged
		}

		return next(ctx)
	}
}
