package bus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getEstimatesForStop(stop int) ([]Bus, error) {
	resp, err := http.Get(fmt.Sprintf("https://emt.kolhos.chichasov.es/api/stop/%v", stop))

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
	resp, err := http.Get("https://emt.kolhos.chichasov.es/api/stop")

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
