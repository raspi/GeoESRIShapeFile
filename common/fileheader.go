package common

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/xerrors"
	"io"
)

const (
	HeaderFileCode         = 9994
	HeaderUnusedMustBe     = 0
	SecondaryHeaderVersion = 1000
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

func (h ShapeFileHeader1) String() string {
	return fmt.Sprintf(`code:%d len:%d %#v`, h.FileCode, h.Length, h.Unused)
}

func (h ShapeFileHeader1) Validate() error {
	if h.FileCode != HeaderFileCode {
		return &InvalidFileCode{Code: h.FileCode}
	}

	for idx, uu := range h.Unused {
		if uu != HeaderUnusedMustBe {
			return &InvalidHeaderUnused{Index: idx, Value: uu}
		}
	}

	if h.Length == 0 {
		return &InvalidHeaderLength{Value: h.Length}
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

	Min struct {
		X, Y float64
	}

	Max struct {
		X, Y float64
	}

	Z struct {
		Min, Max float64
	}

	M struct {
		Min, Max float64
	}
}

func (h ShapeFileHeader2) String() string {
	return fmt.Sprintf(`ver:%d t:%v min:%#v max:%#v Z:%#v M:%#v`, h.Version, h.ShapeType, h.Min, h.Max, h.Z, h.M)
}

func (h ShapeFileHeader2) Validate() error {
	if h.Version != SecondaryHeaderVersion {
		return &InvalidHeaderVersion{Version: h.Version}
	}

	if !IsSupportedShapeType(h.ShapeType) {
		return &ErrInvalidShapeType{ShapeType: h.ShapeType}
	}

	return nil
}

//Read headers shared by .shp and .shx file
func ReadHeaders(r ReadSeekCloser) (err error) {
	err = readFirstHeader(r)
	if err != nil {
		return xerrors.Errorf(`error reading first header (BE) part: %w`, err)
	}

	err = readSecondHeader(r)
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

// Read primary header (notice endianness!)
func readFirstHeader(r ReadSeekCloser) error {
	var hdr1 ShapeFileHeader1
	err := binary.Read(r, binary.BigEndian, &hdr1)
	if err != nil {
		return err
	}

	err = hdr1.Validate()
	if err != nil {
		return err
	}

	return nil
}

// Read secondary header (notice endianness!)
func readSecondHeader(r ReadSeekCloser) error {
	var hdr2 ShapeFileHeader2
	err := binary.Read(r, binary.LittleEndian, &hdr2)
	if err != nil {
		return err
	}

	err = hdr2.Validate()
	if err != nil {
		return err
	}

	return nil
}
