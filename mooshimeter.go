package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type MooshimeterLog struct {
	r        *csv.Reader
	channels []string
}

func NewMooshimeterLog(f io.Reader, channels []string) (m MooshimeterLog) {
	r := csv.NewReader(f)
	r.FieldsPerRecord = len(channels) + 1
	m = MooshimeterLog{
		r:        r,
		channels: channels,
	}
	return
}

func (m MooshimeterLog) Next() (e Entry, err error) {
	record, err := m.r.Read()
	if err != nil {
		if _, ok := err.(*csv.ParseError); ok {
			fmt.Println(err.Error())
			return m.Next()
		}
		return
	}

	values := []float64{}
	for _, v := range record {
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			fmt.Println(err.Error())
			return m.Next()
		}
		values = append(values, f)
	}

	ts := time.Unix(0, int64(values[0]*1000000000.0))
	fields := make(map[string]float64, len(values[1:]))
	for i, v := range values[1:] {
		fields[m.channels[i]] = v
	}

	e = Entry{
		Timestamp: ts,
		Fields:    fields,
	}
	return
}
