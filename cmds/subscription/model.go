package subscription

type channelMessage struct {
	Text string  `json:"text"`
	Link *string `json:"Link"`
}

type channelSubscribers struct {
	ID          string  `bson:"_id"`
	Subscribers []int64 `bson:"subscribers"`
}
