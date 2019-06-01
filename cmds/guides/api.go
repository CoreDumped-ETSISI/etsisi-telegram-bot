package guides

import (
	"encoding/json"
	"net/http"
)

func getAllGuides() (GuideList, error) {
	url := "https://guides.kolhos.chichasov.es/"

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
