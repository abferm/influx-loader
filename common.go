package main

import (
	"time"

	telegraf "github.com/abferm/go-telegraf"
)

type Entry struct {
	Timestamp time.Time
	Fields    map[string]float64
}

func (e *Entry) Shift(shift time.Duration) {
	e.Timestamp = e.Timestamp.Add(shift)
}

func (e Entry) Telegraf(name string, tags map[string]string) telegraf.Measurement {
	m := telegraf.NewMeasurement(name)
	m.SetTime(e.Timestamp)
	m.AddTags(tags)
	for k, v := range e.Fields {
		m.AddFloat64(k, v)
	}
	return m
}

type Log interface {
	Next() (Entry, error)
}
