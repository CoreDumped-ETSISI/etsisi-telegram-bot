package commander

import (
	"fmt"
	"strings"
)

type command struct {
	handler Handler
	prefix  string // /hola
	parsers []argument
}

// Handler is a command handler function that accepts a commander.Context as it's only argument.
type Handler func(Context) error

// Command creates a new command with the specified format. Format admits zero or more parameters
// surrounded with {}. You can specify type with :type. You can also specify ordinality with
// *, + or ? which mean [0, n], [1, n] and [0, 1] respectively.
func (g *CommandGroup) Command(format string, handler Handler) {
	c := &command{
		handler: handler,
		parsers: []argument{},
	}

	parts := strings.Split(format, " ")

	if len(parts) == 0 || parts[0] == "" {
		panic("Format must not be empty!")
	}

	c.prefix = parts[0]

	for _, part := range parts[1:] {
		trm := strings.Trim(part, "{}")
		subp := strings.Split(trm, ":")

		argFlag := argumentSimple
		if strings.HasSuffix(trm, "?") {
			argFlag = argumentOptional
		} else if strings.HasSuffix(trm, "*") {
			argFlag = argumentStar
		} else if strings.HasSuffix(trm, "+") {
			argFlag = argumentPlus
		}

		argType := "string"

		if len(subp) > 1 {
			argType = strings.Trim(subp[1], "?*+")
		}

		parsers := map[string]argumentable{
			"string":   stringArgumentable{},
			"int":      intArgumentable{},
			"float":    floatArgumentable{},
			"duration": durationArgumentable{},
		}

		if _, ok := parsers[argType]; !ok {
			panic(fmt.Sprint("Invalid argument type:", argType))
		}

		c.parsers = append(c.parsers, argument{
			argFlag: argFlag,
			name:    strings.Trim(subp[0], "?*+"),
			parser:  parsers[argType],
		})
	}

	g.commands[c.prefix] = c
}
