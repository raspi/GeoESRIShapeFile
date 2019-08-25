package dbf

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
	"io"
)

var ErrorDeletedRecord = errors.New("deleted record")

type RecordFirstCharacter byte

const (
	DeletedRecord RecordFirstCharacter = 0x2a
	OkRecord      RecordFirstCharacter = 0x20
)

type Record struct {
	Value interface{}
}

func (r Record) String() string {
	return fmt.Sprintf(`'%#v'`, r.Value)
}

func (db *DBaseFile) ReadRecord() (m map[string]Record, err error) {
	if !db.initialized {
		return nil, common.ErrorNotInitialized
	}

	m = make(map[string]Record, db.Header.FieldCount)

	rawalldata := make([]byte, db.Header.RecordSize)
	rBytesAll, err := db.r.Read(rawalldata)
	if err != nil {
		return nil, err
	}

	if rBytesAll != db.Header.RecordSize {
		if rBytesAll == 1 && rawalldata[0] == 0x1a {
			// We are at the end
			return nil, io.EOF
		}

		return nil, fmt.Errorf("full record size mismatch header is %v, had %v:\n%#v", db.Header.RecordSize, rBytesAll, rawalldata[:rBytesAll])
	}

	r := bytes.NewReader(rawalldata)

	first, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch RecordFirstCharacter(first) {
	default:
		return nil, fmt.Errorf(`weird first byte: %[1]d %[1]c %[1]v`, first)
	case DeletedRecord: // deleted record
		return nil, ErrorDeletedRecord
	case OkRecord: // ok
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
