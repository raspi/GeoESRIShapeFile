package dbf

import (
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
	"golang.org/x/xerrors"
	"io"
	"log"
)

type Operation uint8

const (
	KeepAll        Operation = iota // Keep all fields
	KeepOnlyListed                  // Keep only listed fields
	SkipThese                       // Keep everything, except that are listed
)

type DBaseFile struct {
	Header                       Header
	FieldDescriptors             []FieldDescriptor
	converterFunctions           map[string]ConverterFunction
	r                            common.ReadSeekCloser
	debug                        bool
	initialized                  bool
	useDefaultConverterIfMissing bool
	defaultConverter             ConverterFunction
	parseFieldNames              []string
	parseFieldNamesOperation     Operation

	offsets struct {
		mainHeaderEnd int64 // 32
		fieldEnd      int64
		terminatorEnd int64 //
	}
}

func (db *DBaseFile) Close() error {
	return db.r.Close()
}

func (db DBaseFile) Offset() (int64, error) {
	return db.r.Seek(0, io.SeekCurrent)
}

// See datatypes.go for notes about character encoding
func New(fname string, parseFieldNames []string, parseFieldNamesOperation Operation, defaultConverter ConverterFunction, converters map[string]ConverterFunction) (db DBaseFile, err error) {
	f, err := common.OpenFile(fname)
	if err != nil {
		return db, err
	}

	db = DBaseFile{
		converterFunctions:           converters,
		r:                            f,
		debug:                        false,
		useDefaultConverterIfMissing: true,
		defaultConverter:             defaultConverter,
		parseFieldNames:              parseFieldNames,
		parseFieldNamesOperation:     parseFieldNamesOperation,
		initialized:                  false,
	}

	return db, nil
}

func (db *DBaseFile) Initialize() (err error) {
	err = db.readHeader()
	if err != nil {
		return xerrors.Errorf(`error reading header: %w`, err)
	}

	err = db.readFieldHeaders()
	if err != nil {
		return xerrors.Errorf(`error reading field(s): %w`, err)
	}

	err = db.readTerminator()
	if err != nil {
		return xerrors.Errorf(`error reading terminator character after field(s): %w`, err)
	}

	if len(db.FieldDescriptors) != db.Header.FieldCount {
		return fmt.Errorf(`fields found %v but should be %v`, len(db.FieldDescriptors), db.Header.FieldCount)
	}

	if db.debug {
		log.Printf(`header read successfully`)
	}

	db.initialized = true

	return nil
}

func (db *DBaseFile) checkOffset(expected int64, estr string) error {
	currOffset, err := db.Offset()
	if err != nil {
		return err
	}

	if currOffset != expected {
		return fmt.Errorf(`offset is %v, should be %v %v`, currOffset, expected, estr)
	}

	return nil
}

func (db *DBaseFile) SetDebug(flag bool) {
	db.debug = flag
}

func (db *DBaseFile) GetDebug() bool {
	return db.debug
}
