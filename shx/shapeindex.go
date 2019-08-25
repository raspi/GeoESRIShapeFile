package shx

import (
	"encoding/binary"
	"fmt"
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

// Offsets for .shp file
type ShapeIndexRecord struct {
	Offset uint32 // record offset
	Length uint32 // record length
}

func (oi ShapeIndexRecord) String() string {
	return fmt.Sprintf(`offset 0x%04[1]x (%06[1]d) with len 0x%04[2]x (%06[2]d)`, oi.Offset, oi.Length)
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

func (sfi *IndexRecordLookupFile) ReadRecord() (o ShapeIndexRecord, err error) {
	if !sfi.initialized {
		return o, common.ErrorNotInitialized
	}

	err = binary.Read(sfi.r, binary.BigEndian, &o)
	if err != nil {
		return o, err
	}

	o.Offset *= 2
	o.Length *= 2

	sfi.totalFileSize += uint(o.Length)
	sfi.totalFileSize += 8 // Meta data
	sfi.totalRecords++

	return o, nil
}
