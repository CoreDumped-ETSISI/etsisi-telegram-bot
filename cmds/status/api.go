package status

import (
	"encoding/json"
	"net/http"
)

func getStatus() ([]serviceStatus, error) {
	resp, err := http.Get("https://status.kolhos.chichasov.es/api/status")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var services []serviceStatus
	err = json.NewDecoder(resp.Body).Decode(&services)

	if err != nil {
		return nil, err
	}

	return services, nil
}
