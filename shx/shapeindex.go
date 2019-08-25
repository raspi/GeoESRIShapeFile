package shx

import (
	"encoding/binary"
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
	"io"
)

// Offsets for .shp file
type ShapeOffsetIndex struct {
	Offset uint32
	Length uint32
}

func (oi ShapeOffsetIndex) String() string {
	return fmt.Sprintf(`offset 0x%04[1]x (%06[1]d) with len 0x%04[2]x (%06[2]d)`, oi.Offset, oi.Length)
}

func New(fname string) (offsets []ShapeOffsetIndex, shpTotalFilesize int64, err error) {
	f, err := common.OpenFile(fname)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	err = common.ReadHeaders(f)
	if err != nil {
		return nil, 0, err
	}

	totalLen := int64(100) // header = 100

	for {
		var o ShapeOffsetIndex
		err = binary.Read(f, binary.BigEndian, &o)
		if err == io.EOF {
			break
		}

		if err != nil {
			return offsets, totalLen, err
		}

		o.Offset *= 2
		o.Length *= 2

		totalLen += int64(o.Length) + 8 // metadata
		offsets = append(offsets, o)
	}

	return offsets, totalLen, nil
}
