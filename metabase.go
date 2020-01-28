package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type MetaBaseLog struct {
	r          *csv.Reader
	fieldNames []string
}

func NewMetaBaseLog(f io.Reader) (l MetaBaseLog, err error) {
	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return
	}
	fieldNames := []string{}
	for _, name := range header[3:] {
		name = strings.Replace(name, " ", "_", -1)
		name = strings.Replace(name, "(", "", -1)
		name = strings.Replace(name, ")", "", -1)
		fieldNames = append(fieldNames, name)
	}
	l = MetaBaseLog{
		r:          r,
		fieldNames: fieldNames,
	}
	return
}

func (m MetaBaseLog) Next() (e Entry, err error) {
	record, err := m.r.Read()
	if err != nil {
		if _, ok := err.(*csv.ParseError); ok {
			fmt.Println(err.Error())
			return m.Next()
		}
		return
	}

	tsMillis, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return m.Next()
	}

	values := []float64{}
	for _, v := range record[3:] {
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			fmt.Println(err.Error())
			return m.Next()
		}
		values = append(values, f)
	}

	ts := time.Unix(0, tsMillis*1000000)
	fields := make(map[string]float64, len(values))
	for i, v := range values {
		fields[m.fieldNames[i]] = v
	}

	e = Entry{
		Timestamp: ts,
		Fields:    fields,
	}
	return
}
