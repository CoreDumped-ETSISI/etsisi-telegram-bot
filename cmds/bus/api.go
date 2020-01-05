package bus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

func getEstimatesForStop(stop int) ([]Bus, error) {
	resp, err := http.Get(fmt.Sprintf(services.Get("emt", 8080) + "/api/stop/%v", stop))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var arrive []Bus
	err = json.NewDecoder(resp.Body).Decode(&arrive)

	if err != nil {
		return nil, err
	}

	return arrive, nil
}

func getUniEstimates() (*UniversityStops, error) {
	resp, err := http.Get(services.Get("emt", 8080) + "/api/stop")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var arrive UniversityStops
	err = json.NewDecoder(resp.Body).Decode(&arrive)

	if err != nil {
		return nil, err
	}

	return &arrive, nil
}
