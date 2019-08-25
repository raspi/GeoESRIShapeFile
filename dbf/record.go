package dbf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var ErrorDeletedRecord = errors.New("deleted record")

type Record struct {
	Value interface{}
}

func (r Record) String() string {
	return fmt.Sprintf(`'%#v'`, r.Value)
}

func (db *DBaseFile) ReadRecord() (m map[string]Record, err error) {
	m = make(map[string]Record, len(db.FieldDescriptors))

	recordSizeLen := db.Header.RecordSize

	rawalldata := make([]byte, recordSizeLen)
	rBytesAll, err := db.r.Read(rawalldata)
	if err != nil {
		return nil, err
	}

	if rBytesAll != recordSizeLen {
		if rBytesAll == 1 && rawalldata[0] == 0x1a {
			// We are at the end
			return nil, io.EOF
		}

		return nil, fmt.Errorf("full record size mismatch header is %v, had %v:\n%#v", recordSizeLen, rBytesAll, rawalldata[:rBytesAll])
	}

	r := bytes.NewReader(rawalldata)

	first, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch first {
	default:
		return nil, fmt.Errorf(`weird first byte: %[1]d %[1]c %[1]v`, first)
	case 0x2a: // deleted record
		return nil, ErrorDeletedRecord
	case 0x20: // ok
	}

	for _, f := range db.FieldDescriptors {

		if !isSupportedDataType(f.Type) {
			return nil, NewErrorNotSupportedDataType(f.Type)
		}

		skip := false // default keep all

		if db.parseFieldNamesOperation != KeepAll {
			if db.parseFieldNamesOperation == KeepOnlyListed {
				// skip all by default
				skip = true
			}

			for _, fn := range db.parseFieldNames {
				if fn == f.Name {
					switch db.parseFieldNamesOperation {
					case KeepOnlyListed:
						skip = false
					case SkipThese:
						skip = true
					}

					break
				}
			}
		}

		if skip {
			continue
		}

		rawdata := make([]byte, f.Length)
		rBytes, err := r.Read(rawdata)
		if err != nil {
			return nil, err
		}

		if rBytes != f.Length {
			return nil, fmt.Errorf(`record size mismatch %v != %v`, rBytes, f.Length)
		}

		// Find converter
		converter, ok := db.converterFunctions[f.Name]

		if !ok && !db.useDefaultConverterIfMissing {
			return nil, fmt.Errorf(`no such converter: %v %v`, f.Type, f.Name)
		}

		if !ok {
			converter = db.defaultConverter
		}

		rec, err := converter(rawdata)
		if err != nil {
			return nil, err
		}

		m[f.Name] = rec

	}

	return m, nil
}
