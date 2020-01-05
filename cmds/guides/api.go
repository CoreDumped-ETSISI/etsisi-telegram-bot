package guides

import (
	"encoding/json"
	"net/http"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

func getAllGuides() (GuideList, error) {
	url := services.Get("guias", 80)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var gl GuideList

	err = json.NewDecoder(resp.Body).Decode(&gl)

	if err != nil {
		return nil, err
	}

	return gl, nil
}
