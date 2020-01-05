package services

import (
	"os"
	"strings"
	"fmt"
)

func Get(name string, defaultPort int) string {
	if env, ok := os.LookupEnv("HOST_" + strings.ToUpper(name)); ok {
		return env
	}

	if defaultPort > 0 {
		return fmt.Sprintf("http://%v:%v", name, defaultPort)
	}
	return fmt.Sprintf("http://%s", name)
}