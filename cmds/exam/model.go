package exam

import "time"

type Exam struct {
	Name     string    `json:"name"`
	RawName  string    `json:"name_raw"`
	Aulas    []string  `json:"aulas"`
	Timeslot string    `json:"timeslot"`
	Day      string    `json:"day"`
	Date     string    `json:"date"`
	DateTime time.Time `json:"-"`
	Tags     []string  `json:"tags"`
}

var (
	Grados = map[string]string{
		"Software":     "software",
		"Computadores": "compu",
		"SI":           "si",
		"TI":           "ti",
		"Optativas":    "optativa",
	}
)
