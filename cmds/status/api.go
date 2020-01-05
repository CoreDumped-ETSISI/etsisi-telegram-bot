package status

import (
	"encoding/json"
	"net/http"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

func getStatus() ([]serviceStatus, error) {
	resp, err := http.Get(services.Get("unistatus", 8889)+"/api/status")

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
