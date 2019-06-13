package subscription

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DBContext struct {
	DB *mgo.Session
}

func New(addr, db, username, pass string) *DBContext {
	info := &mgo.DialInfo{
		Addrs:    []string{addr},
		Database: db,
		Username: username,
		Password: pass,
		Timeout:  10 * time.Second,
	}

	database, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	return &DBContext{
		DB: database,
	}
}

func (s *DBContext) getSubscribers(channel string) ([]int64, error) {
	sesh := s.DB.Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("subscriptions")

	var ch channelSubscribers
	err := col.FindId(channel).One(&ch)

	if err != nil {
		if err == mgo.ErrNotFound {
			return []int64{}, nil
		}
		return nil, err
	}

	return ch.Subscribers, nil
}

func (s *DBContext) addSubscriber(channel string, user int64) error {
	sesh := s.DB.Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("subscriptions")

	err := col.UpdateId(channel, bson.M{
		"$addToSet": bson.M{
			"subscribers": user,
		},
	})

	if err != nil {
		if err == mgo.ErrNotFound {
			var ch channelSubscribers
			ch.ID = channel
			ch.Subscribers = []int64{user}
			return col.Insert(ch)
		}

		return err
	}

	return nil
}

func (s *DBContext) removeSubscriber(channel string, user int64) error {
	sesh := s.DB.Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("subscriptions")

	err := col.UpdateId(channel, bson.M{
		"$pull": bson.M{
			"subscribers": user,
		},
	})

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}

		return err
	}

	return nil
}
