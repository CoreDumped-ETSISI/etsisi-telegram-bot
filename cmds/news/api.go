package news

import (
	"encoding/json"
	"net/http"
	"time"
)

func fetchFeed(feed string) ([]newsItem, error) {
	url := "https://uninews.kolhos.chichasov.es/" + feed

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
