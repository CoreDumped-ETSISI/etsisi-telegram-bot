package main

import (
	"strings"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

func filterEmpty(s []string) []string {
	out := make([]string, 0, len(s))

	for i := range s {
		if s[i] != "" {
			out = append(out, s[i])
		}
	}

	return out
}

func getChatTitle(msg *tb.Message) string {
	chatTitle := msg.Chat.Title

	if chatTitle == "" {
		chatTitle = strings.Join(filterEmpty(
			[]string{
				msg.Chat.FirstName,
				msg.Chat.LastName,
			}), " ")
	}

	return chatTitle
}

func getSenderName(sender *tb.User) string {
	parts := []string{
		sender.FirstName,
		sender.LastName,
	}

	parts = filterEmpty(parts)

	if len(parts) == 0 {
		return sender.UserName
	}

	if sender.UserName != "" {
		parts = append(parts, "(" + sender.UserName + ")")
	}

	return strings.Join(parts, " ")
}
