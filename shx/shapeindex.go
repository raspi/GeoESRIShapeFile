package shx

import (
	"encoding/binary"
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
	"io"
	"log"
)

type ShapeFileIndex struct {
	r             common.ReadSeekCloser
	debug         bool
	initialized   bool
	totalFileSize uint
	totalRecords  uint
}

// Offsets for .shp file
type ShapeOffsetIndex struct {
	Offset uint32 // record offset
	Length uint32 // record length
}

func (oi ShapeOffsetIndex) String() string {
	return fmt.Sprintf(`offset 0x%04[1]x (%06[1]d) with len 0x%04[2]x (%06[2]d)`, oi.Offset, oi.Length)
}

func New(fname string) (sfi ShapeFileIndex, err error) {
	f, err := common.OpenFile(fname)
	if err != nil {
		return sfi, err
	}

	return ShapeFileIndex{
		r:           f,
		debug:       false,
		initialized: false,
	}, nil
}

func (sfi *ShapeFileIndex) SetDebug(flag bool) {
	sfi.debug = flag
}

func (sfi *ShapeFileIndex) GetDebug() bool {
	return sfi.debug
}

func (sfi *ShapeFileIndex) Close() error {
	return sfi.r.Close()
}

func (sfi *ShapeFileIndex) Initialize() (err error) {
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

func (sfi ShapeFileIndex) GetShapeFileSize() uint {
	return sfi.totalFileSize
}

func (sfi ShapeFileIndex) GetRecordCount() uint {
	return sfi.totalRecords
}

func (sfi *ShapeFileIndex) ReadRecord() (o ShapeOffsetIndex, err error) {
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
