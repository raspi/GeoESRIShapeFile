package dbf

import (
	"strconv"
	"strings"
)

var DefaultConverterToInt = func(data []byte) (rec Record, err error) {
	s := strings.Trim(string(data), ` `)

	if s == `` {
		rec.Value = nil
	} else {
		rec.Value, err = strconv.ParseInt(s, 10, 64)
	}

	return rec, err
}

var DefaultConverterToString = func(data []byte) (rec Record, err error) {
	s := strings.Trim(string(data), ` `)

	if s == `` {
		rec.Value = nil
	} else {
		rec.Value = s
	}

	return rec, nil
}

var DefaultConverterToError = func(data []byte) (rec Record, err error) {
	rec.Value = data
	return rec, ErrorConverterNotFound
}
