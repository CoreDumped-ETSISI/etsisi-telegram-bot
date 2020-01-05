package horario

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/services"
)

var (
	NoSuchGroupError = errors.New("Ese grupo no existe!")
)

func getAllHorarios() (horarios, error) {
	resp, err := http.Get(services.Get("horarios", 80)+"/horario")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data horarios
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getHorarioForGroup(group string) ([][]class, error) {
	resp, err := http.Get(services.Get("horarios", 80)+"/horario/" + group)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		// Group not found
		return nil, NoSuchGroupError
	}

	var data [][]class
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
