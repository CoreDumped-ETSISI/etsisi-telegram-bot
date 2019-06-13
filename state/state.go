package state

import (
	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
)

type T interface {
	Redis() *redis.Client
	Mongo() *mgo.Session
}
