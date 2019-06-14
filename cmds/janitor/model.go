package janitor

type banEvent struct {
	UserID  int   `json:"user_id,omitempty"`
	AdminID int   `json:"admin_id,omitempty"`
	ChatID  int64 `json:"chat_id,omitempty"`
}

type channelManagement struct {
	ChatID   int64  `bson:"_id,omitempty"`
	AdminsID []int  `bson:"admins_id,omitempty"`
	Public   bool   `bson:"public"`
	Name     string `bson:"name,omitempty"`
}
