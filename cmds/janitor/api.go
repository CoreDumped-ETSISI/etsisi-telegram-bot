package janitor

import "github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

func isChatManaged(state state.T, chatid int64) (bool, error) {
	sesh := state.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("managed-channels")

	n, err := col.FindId(chatid).Count()

	if err != nil {
		return false, err
	}

	return n > 0, nil
}
