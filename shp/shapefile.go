package shp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/raspi/GeoESRIShapeFile/common"
	"io"
	"log"
)

type RecordHeader struct {
	Number uint32
	Length uint32
}

type ShapeFile struct {
	r           common.ReadSeekCloser
	debug       bool
	initialized bool
}

func (sf *ShapeFile) Close() error {
	return sf.r.Close()
}

func (sf *ShapeFile) SetDebug(flag bool) {
	sf.debug = flag
}

func (sf ShapeFile) GetDebug() bool {
	return sf.debug
}

func (sf *ShapeFile) ReadRecordAt(offset int64) (idx uint32, record ShapeTypeI, err error) {
	_, err = sf.r.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, nil, err
	}

	return sf.ReadRecord()
}

func (sf *ShapeFile) ReadRecord() (idx uint32, record ShapeTypeI, err error) {
	if !sf.initialized {
		return idx, nil, common.ErrorNotInitialized
	}

	offset := int64(-1)
	if sf.debug {
		offset, err = sf.r.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, nil, err
		}
	}

	var rechdr RecordHeader

	err = binary.Read(sf.r, binary.BigEndian, &rechdr)
	if err != nil {
		return 0, nil, err
	}

	rechdr.Length *= 2
	rechdr.Number--

	rawshapedata := make([]byte, rechdr.Length)
	rBytes, err := sf.r.Read(rawshapedata)
	if err != nil {
		return 0, nil, err
	}

	if uint32(rBytes) != rechdr.Length {
		return 0, nil, fmt.Errorf(`read %v but len is %v?`, rBytes, rechdr.Length)
	}

	if sf.debug {
		log.Printf(`Read shape #%[1]v with len 0x%04[2]x (%06[2]d) at offset 0x%04[3]x (%06[3]d)`, rechdr.Number, rechdr.Length, offset)
	}

	rec, err := sf.readRecordData(bytes.NewReader(rawshapedata))
	if err != nil {
		return 0, nil, err
	}

	err = rec.Validate()
	if err != nil {
		return 0, nil, err
	}

	return rechdr.Number, rec, nil
}

func (sf *ShapeFile) readRecordData(r io.ReadSeeker) (rec ShapeTypeI, err error) {
	var shapeType common.ShapeType
	err = binary.Read(r, binary.LittleEndian, &shapeType)
	if err != nil {
		return nil, err
	}

	//log.Printf(`got shape %v`, shapeType)

	switch shapeType {
	case common.POLYLINEZ:
		var nrec PolyLineZ
		return nrec.read(r)
	default:
		return nil, fmt.Errorf(`unknown shape style: %v`, shapeType)
	}
}

func New(fname string) (sf ShapeFile, err error) {
	sf.debug = false

	f, err := common.OpenFile(fname)
	if err != nil {
		return sf, err
	}

	return ShapeFile{
		r:           f,
		initialized: false,
	}, nil
}

func (sf *ShapeFile) Initialize() (err error) {
	err = common.ReadHeaders(sf.r)
	if err != nil {
		return err
	}

	if sf.debug {
		log.Printf(`header read successfully`)
	}

	sf.initialized = true

	return nil
}
