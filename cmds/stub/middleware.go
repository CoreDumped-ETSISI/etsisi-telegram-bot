package stub

import (
	"github.com/guad/commander"
)

func Middleware(next commander.Handler) commander.Handler {
	return func(ctx commander.Context) error {
		return Cmd(ctx)
	}
}
