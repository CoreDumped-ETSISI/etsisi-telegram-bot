package salas

import (
	"encoding/json"
	"net/http"
)

type salasResponse struct {
	Salas []salaTimesheet `json:"salas"`
}

type salaTimesheet struct {
	ID       int        `json:"id"`
	Occupied []timeSlot `json:"occupied"`
}

type timeSlot struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func getSalas() (*salasResponse, error) {
	resp, err := http.Get("https://biblio.kolhos.chichasov.es/api/salas")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var salas salasResponse
	err = json.NewDecoder(resp.Body).Decode(&salas)

	if err != nil {
		return nil, err
	}

	return &salas, nil
}
