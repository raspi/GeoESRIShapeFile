package common

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/xerrors"
	"io"
)

/*
Headers shared by .shp and .shx
*/

/*
Byte 0  File Code    9994        Integer Big
Byte 4  Unused       0           Integer Big
Byte 8  Unused       0           Integer Big
Byte 12 Unused       0           Integer Big
Byte 16 Unused       0           Integer Big
Byte 20 Unused       0           Integer Big
Byte 24 File Length  File Length Integer Big
*/
type ShapeFileHeader1 struct {
	FileCode uint32    // 9994
	Unused   [5]uint32 // all 0
	Length   uint32
}

func (header1 ShapeFileHeader1) Validate() error {
	if header1.FileCode != 9994 {
		return fmt.Errorf(`filecode was not 9994`)
	}

	for _, uu := range header1.Unused {
		if uu != 0 {
			return fmt.Errorf(`unused was not 0`)
		}
	}

	return nil
}

/*
Byte 28  Version  1000            Integer Little
Byte 32  Shape    Type Shape Type Integer Little
Byte 36  Bounding Box Xmin        Double  Little
Byte 44  Bounding Box Ymin        Double  Little
Byte 52  Bounding Box Xmax        Double  Little
Byte 60  Bounding Box Ymax        Double  Little
Byte 68* Bounding Box Zmin        Double  Little
Byte 76* Bounding Box Zmax        Double  Little
Byte 84* Bounding Box Mmin        Double  Little
Byte 92* Bounding Box Mmax        Double  Little
*/
type ShapeFileHeader2 struct {
	Version   uint32
	ShapeType ShapeType
	BBoxXMin  float64
	BBoxYMin  float64
	BBoxXMax  float64
	BBoxYMax  float64
	BBoxZMin  float64
	BBoxZMax  float64
	BBoxMMin  float64
	BBoxMMax  float64
}

//Read headers shared by .shp and .shx file
func ReadHeaders(r ReadSeekCloser) error {
	// Read primary header (notice endianness!)
	var hdr1 ShapeFileHeader1
	err := binary.Read(r, binary.BigEndian, &hdr1)
	if err != nil {
		return xerrors.Errorf(`error reading first header (BE) part: %w`, err)
	}

	err = hdr1.Validate()
	if err != nil {
		return err
	}

	// Read secondary header (notice endianness!)
	var hdr2 ShapeFileHeader2
	err = binary.Read(r, binary.LittleEndian, &hdr2)
	if err != nil {
		return xerrors.Errorf(`error reading second header (LE) part: %w`, err)
	}

	offset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	if offset != 100 {
		return fmt.Errorf(`offset is not 100`)
	}

	return nil
}
