package state

import (
	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

type T interface {
	Redis() *redis.Client
	Mongo() *mgo.Session
	Bot() *tb.BotAPI
}
