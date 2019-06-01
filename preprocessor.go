package main

import (
	"strings"

	"github.com/guad/commander"
)

type CustomTelegramPreprocessor struct {
	BotName string
}

func (t *CustomTelegramPreprocessor) Process(args []string) ([]string, bool) {
	orig := &commander.TelegramPreprocessor{
		BotName: t.BotName,
	}

	if newargs, ok := orig.Process(args); !ok {
		return nil, false
	} else {
		args = newargs
	}

	// Convert commands like /mycmd_Arg1_Arg2 into /mycmd Arg1 Arg2
	if len(args) == 1 && strings.ContainsRune(args[0], '_') {
		args = strings.Split(args[0], "_")
	}

	return args, true
}
