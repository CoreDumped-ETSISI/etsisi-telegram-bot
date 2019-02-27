package main

import (
	"fmt"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getChatTitle(msg *tb.Message) string {
	chatTitle := msg.Chat.Title

	if chatTitle == "" {
		chatTitle = fmt.Sprint(msg.Chat.FirstName, msg.Chat.LastName)
	}

	return chatTitle
}

func getSenderName(sender *tb.User) string {
	if sender.FirstName == "" && sender.LastName == "" {
		return sender.UserName
	}

	name := sender.FirstName

	if sender.FirstName == "" {
		name = sender.LastName
	} else if sender.LastName != "" {
		name += " " + sender.LastName
	}

	if sender.UserName != "" {
		name = fmt.Sprintf("%v (%v)", name, sender.UserName)
	}

	return name
}
