package shx

import (
	"encoding/binary"
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
)

// Offsets for .shp file
type ShapeIndexRecord struct {
	Offset uint32 // record offset
	Length uint32 // record length
}

func (oi ShapeIndexRecord) String() string {
	return fmt.Sprintf(`offset 0x%04[1]x (%06[1]d) with len 0x%04[2]x (%06[2]d)`, oi.Offset, oi.Length)
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
