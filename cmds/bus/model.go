package bus

import (
	"time"
)

type Bus struct {
	LineID      string        `json:"lineId"`
	Destination string        `json:"destination"`
	BusID       string        `json:"busId"`
	TimeLeft    time.Duration `json:"timeLeft"`
	Distance    int           `json:"distance"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
}

type UniversityStops struct {
	SentidoSierra []Bus `json:"sierra"`
	SentidoConde  []Bus `json:"conde"`
}
