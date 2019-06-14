package janitor

import (
	"errors"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	ErrNotSupergroup = errors.New("This group must be a supergroup")
)

func GetInviteLink(bot *tb.BotAPI, chatid int64) (string, error) {
	link, err := bot.GetChat(tb.ChatConfig{ChatID: chatid})

	if err != nil {
		return "", err
	}

	if link.InviteLink == "" {
		return "", ErrNotSupergroup
	}

	return link.InviteLink, nil
}
