package main

import (
	"strings"
	"testing"
)

func TestPreprocessor(t *testing.T) {
	p := CustomTelegramPreprocessor{
		BotName: "testan",
	}

	var cases = []struct {
		command string
		expect  []string
	}{
		{"hello@testan", []string{"hello"}},
		{"hello", []string{"hello"}},
		{"hello_world", []string{"hello", "world"}},
		{"hello_world_hello", []string{"hello", "world", "hello"}},
		{"hello_world whats up", []string{"hello_world", "whats", "up"}},
	}

	for _, data := range cases {
		args, ok := p.Process(strings.Split(data.command, " "))

		if !ok {
			t.Fail()
			t.Logf("Unexpected failure for input %#v", data.command)
		} else if len(args) != len(data.expect) {
			t.Fail()
			t.Logf("Expected %+v, got %+v from %#v", data.expect, args, data.command)
		} else {
			for i := range data.expect {
				if data.expect[i] != args[i] {
					t.Fail()
					t.Logf("Expected %+v, got %+v from %#v", data.expect, args, data.command)
					return
				}
			}
		}
	}
}
