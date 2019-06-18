package verify

import (
	"math/rand"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"
)

func IsUserVerified(userid int) bool {
	sesh := state.G.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("verified-users")

	var vu verifiedUser
	err := col.FindId(userid).One(&vu)

	return err == nil
}

func startNewVerification(userid int) (string, error) {
	// Generate random session ID
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	n := 32
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	token := string(b)

	err := state.G.Redis().Set("VERIFY_SESS_"+token, userid, 30*time.Minute).Err()

	if err != nil {
		return "", nil
	}

	return token, nil
}

func buildVerificationURL(token string) string {
	// TODO
	return ""
}

func verifyUser(userid int) error {
	sesh := state.G.Mongo().Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("verified-users")

	vu := verifiedUser{
		UserID: userid,
		Date:   time.Now(),
	}

	return col.Insert(vu)
}
