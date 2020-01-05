package news

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

func fetchFeed(feed string) ([]newsItem, error) {
	url := services.Get("uninews", 80) + "/" + feed

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var news []newsItem
	err = json.NewDecoder(resp.Body).Decode(&news)

	if err != nil {
		return nil, err
	}

	return news, nil
}
