package main

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

type config struct {
	commandTimeout time.Duration
	redis          *redis.Client
	logLevel       int
	bot            *tb.BotAPI
}

func newConfig() config {
	cfg := config{}

	timeout, _ := time.ParseDuration(os.Getenv("MESSAGE_TIMEOUT"))

	cfg.commandTimeout = timeout

	redisb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg.redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       redisb,
	})

	cfg.logLevel = 4 // INFO

	lvlEnv := os.Getenv("LOG_LEVEL")

	/*
		PanicLevel = 0,
		FatalLevel = 1,
		ErrorLevel = 2,
		WarnLevel  = 3,
		InfoLevel  = 4,
		DebugLevel = 5,
		TraceLevel = 6,
	*/

	if lvlEnv != "" {
		cfg.logLevel, _ = strconv.Atoi(lvlEnv)
	}

	return cfg
}