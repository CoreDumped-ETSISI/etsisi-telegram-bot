package verify

import "time"

type userVerified struct {
	UserID int `json:"user_id,omitempty"`
}

type verifiedUser struct {
	UserID int       `json:"_id,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}
