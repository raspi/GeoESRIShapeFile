package shx

import (
	"github.com/raspi/GeoESRIShapeFile/common"
	"io"
	"log"
)

type IndexRecordLookupFile struct {
	r             common.ReadSeekCloser
	debug         bool
	initialized   bool
	totalFileSize uint
	totalRecords  uint
}

func New(fname string) (sfi IndexRecordLookupFile, err error) {
	f, err := common.OpenFile(fname)
	if err != nil {
		return sfi, err
	}

	return IndexRecordLookupFile{
		r:           f,
		debug:       false,
		initialized: false,
	}, nil
}

func (sfi *IndexRecordLookupFile) SetDebug(flag bool) {
	sfi.debug = flag
}

func (sfi *IndexRecordLookupFile) GetDebug() bool {
	return sfi.debug
}

func (sfi *IndexRecordLookupFile) Close() error {
	return sfi.r.Close()
}

func (sfi *IndexRecordLookupFile) Initialize() (err error) {
	err = common.ReadHeaders(sfi.r)
	if err != nil {
		return err
	}

	offset, err := sfi.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	sfi.totalFileSize += uint(offset)

	if sfi.debug {
		log.Printf(`header read successfully`)
	}

	sfi.initialized = true
	return nil
}

func (sfi IndexRecordLookupFile) GetShapeFileSize() uint {
	return sfi.totalFileSize
}

func (sfi IndexRecordLookupFile) GetRecordCount() uint {
	return sfi.totalRecords
}
