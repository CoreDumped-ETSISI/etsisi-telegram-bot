package main

import (
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
)

func TestChatTitle(t *testing.T) {
	var cases = []struct {
		title     string
		firstName string
		lastName  string
		expected  string
	}{
		{"Grupo uno", "", "", "Grupo uno"},
		{"Grupo dos", "Grupo", "last", "Grupo dos"},
		{"", "Felipe", "Perez", "Felipe Perez"},
		{"", "Felipe", "", "Felipe"},
		{"", "", "Perez", "Perez"},
		{"", "", "", ""},
	}

	for _, data := range cases {
		msg := tb.Message{
			Chat: &tb.Chat{
				Title:     data.title,
				FirstName: data.firstName,
				LastName:  data.lastName,
			},
		}

		title := getChatTitle(&msg)

		if title != data.expected {
			t.Fail()
			t.Logf("Expected: %#v, got %#v from %+v", data.expected, title, data)
		}
	}
}

func TestSenderName(t *testing.T) {
	var cases = []struct {
		username  string
		firstName string
		lastName  string
		expected  string
	}{
		{"guad", "", "", "guad"},
		{"guad", "Felipe", "Perez", "Felipe Perez (guad)"},
		{"", "Felipe", "Perez", "Felipe Perez"},
		{"", "Felipe", "", "Felipe"},
		{"", "", "Perez", "Perez"},
		{"", "", "", ""},
		{"guad", "Felipe", "", "Felipe (guad)"},
		{"guad", "", "Perez", "Perez (guad)"},
	}

	for _, data := range cases {
		user := tb.User{
			UserName:  data.username,
			FirstName: data.firstName,
			LastName:  data.lastName,
		}

		name := getSenderName(&user)

		if name != data.expected {
			t.Fail()
			t.Logf("Expected: %#v, got %#v from %+v", data.expected, name, data)
		}
	}
}
