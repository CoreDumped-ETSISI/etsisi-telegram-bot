package stub

import (
	"encoding/json"
	"github.com/guad/commander"
	"io/ioutil"
	"os"
	log "github.com/sirupsen/logrus"
	"strings"
)

func emptyMiddleware(next commander.Handler) commander.Handler {
	return next
}

func logError(err error) {
	log.WithError(err).Error("Failed to start dynamic stubbing")
}

func Middleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		return Cmd(ctx, "Este comando ha sido temporalmente desactivado.")
	}
}

func DynamicMiddleware() commander.CommandMiddleware {
	configPath := "/disabledCommands.json"

	if e, ok := os.LookupEnv("STUB_CMDS_CONFIG"); ok {
		configPath = e
	}

	f, err := os.Open(configPath)

	if err != nil {
		logError(err)
		return emptyMiddleware
	}


	b, err := ioutil.ReadAll(f)

	if err != nil {
		logError(err)
		return emptyMiddleware
	}

	var disabled map[string]struct{
		Message string
		Disabled bool
	}

	err = json.Unmarshal(b, &disabled)

	if err != nil {
		logError(err)
		return emptyMiddleware
	}

	keys := []string{}

	for k := range disabled {
		keys = append(keys, k)
	}

	log.WithField("commands", strings.Join(keys, ", ")).Debug("Stubbed out commands")

	return func(next commander.Handler) commander.Handler {
		return func(ctx commander.Context) error {
			log.WithField("command", ctx.Name).Debug("Trying to stub out command")

			if data, ok := disabled[ctx.Name]; ok && data.Disabled {
				log.WithField("command", ctx.Name).Debug("Stubbed out command!")
				return Cmd(ctx, data.Message)
			}
			return next(ctx)
		}
	}
}
