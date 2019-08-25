package dbf

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// dBase Field Descriptor header after main header
// size: 32 bytes
// count: N
type rawFieldDescriptor struct {
	Name           [11]byte // Field name
	Type           DataType // Field type
	_              [4]byte
	Length         uint8 // Field length
	DecimalCount   uint8 // For floats
	_              [2]byte
	WorkAreaID     byte
	_              [2]byte
	FlagSetField   byte
	_              [7]byte
	IndexFieldFlag byte
}

func (r rawFieldDescriptor) String() string {
	return fmt.Sprintf(`%c len:%03d dc:%03v %s`, r.Type, r.Length, r.DecimalCount, r.Name)
}

// FieldDescriptor describes a field. For example: Street address which is Character and 200 bytes long
type FieldDescriptor struct {
	Name           string   // Field name "STREETNAME"
	Type           DataType // Field type, like a Character, Nemerical, etc
	Length         int
	DecimalCount   int // For floats
	WorkAreaID     uint8
	MdxFlag        uint8
	FlagSetField   byte
	IndexFieldFlag byte
}

func (fd FieldDescriptor) String() string {
	return fmt.Sprintf(`%s %v len:%v deccount:%v waID:%v mdx:%v`, fd.Name, fd.Type, fd.Length, fd.DecimalCount, fd.WorkAreaID, fd.WorkAreaID)
}

// Read raw Fields for descriptions after main header
func (db *DBaseFile) readFieldHeaders() (err error) {
	err = db.checkOffset(db.offsets.mainHeaderEnd, `starting to read field descriptions`)
	if err != nil {
		return err
	}

	rawf := make([]rawFieldDescriptor, db.Header.FieldCount)
	err = binary.Read(db.r, binary.LittleEndian, &rawf)
	if err != nil {
		return err
	}

	err = db.checkOffset(db.offsets.fieldEnd, `after reading field descriptions`)
	if err != nil {
		return err
	}

	// Build proper fields
	for _, f := range rawf {
		db.FieldDescriptors = append(db.FieldDescriptors, FieldDescriptor{
			Name:         strings.TrimRight(string(f.Name[:]), "\x00"),
			Type:         f.Type,
			Length:       int(f.Length),
			DecimalCount: int(f.DecimalCount), // for floats
			WorkAreaID:   f.WorkAreaID,
			//MdxFlag:      f.MdxFlag,
		})
	}

	return nil
}
