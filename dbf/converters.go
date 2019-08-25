package dbf

import "errors"

var ErrorConverterNotFound = errors.New("no converter for this type")

type ConverterFunction func(data []byte) (rec Record, err error)
