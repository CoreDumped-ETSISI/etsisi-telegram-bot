package commander

import (
	"strconv"
	"time"
)

const (
	argumentSimple = iota
	argumentOptional
	argumentStar
	argumentPlus
)

type argument struct {
	parser  argumentable
	name    string
	argFlag int
}

type argumentable interface {
	Parse(text string) (interface{}, error)
}

type stringArgumentable struct{}

func (s stringArgumentable) Parse(text string) (interface{}, error) {
	return text, nil
}

type intArgumentable struct{}

func (s intArgumentable) Parse(text string) (interface{}, error) {
	return strconv.Atoi(text)
}

type floatArgumentable struct{}

func (s floatArgumentable) Parse(text string) (interface{}, error) {
	return strconv.ParseFloat(text, 64)
}

type durationArgumentable struct{}

func (s durationArgumentable) Parse(text string) (interface{}, error) {
	return time.ParseDuration(text)
}
