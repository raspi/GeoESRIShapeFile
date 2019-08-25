package dbf

import (
	"encoding/binary"
	"fmt"
	"time"
)

// raw dBase header
// size: 32 bytes
// count: 1
type rawHeader struct {
	Version                              Version // 1
	UpdateYear, UpdateMonth, UpdateDay   uint8   // 2
	RecordCount                          uint32  // 8
	LengthHeaderBytes, LengthRecordBytes uint16  // 10
	_                                    [2]byte // 14
	IncompleteTransactionFlag            uint8   // 15
	EncryptionFlag                       uint8   // 16
	FreeRecThread                        uint32
	_                                    [8]byte
	MdxFlag                              uint8   // 29
	LanguageDriver                       uint8   // 30
	_                                    [2]byte // 32
}

func (rh rawHeader) String() string {
	return fmt.Sprintf(`ver:%v updated: %d-%d-%d records:%d hdrsize:%d recsize:%d`, rh.Version, rh.UpdateYear, rh.UpdateMonth, rh.UpdateDay, rh.RecordCount, rh.LengthHeaderBytes, rh.LengthRecordBytes)
}

// Proper header for this library
type Header struct {
	Version Version
	Date    time.Time

	FieldCount int // How many fields

	RecordCount int // How many records?
	RecordSize  int // How many bytes is each record
}

// Read main header
func (db *DBaseFile) readHeader() (err error) {
	err = db.checkOffset(0, `before reading main header`)
	if err != nil {
		return err
	}

	var rawhdr rawHeader
	rawHeaderBinSize := int64(binary.Size(rawhdr))
	if rawHeaderBinSize != 32 {
		return fmt.Errorf(`raw header size is %v, should be 32`, rawHeaderBinSize)
	}

	err = binary.Read(db.r, binary.LittleEndian, &rawhdr)
	if err != nil {
		return err
	}

	offset, err := db.Offset()

	err = db.checkOffset(offset, `after reading main header`)
	if err != nil {
		return err
	}

	// Pre-count offsets and field counts

	rawFieldBinSize := binary.Size(rawFieldDescriptor{})
	if rawFieldBinSize != 32 {
		return fmt.Errorf(`rawField size should be 32, is %v`, rawFieldBinSize)
	}

	rawFieldCount := int(rawhdr.LengthHeaderBytes)/rawFieldBinSize - 1

	db.offsets.mainHeaderEnd = offset
	db.offsets.fieldEnd = int64(rawhdr.LengthHeaderBytes) - 1
	db.offsets.terminatorEnd = int64(rawhdr.LengthHeaderBytes)

	// Build proper header
	db.Header = Header{
		Version:     rawhdr.Version,
		Date:        time.Date(int(rawhdr.UpdateYear)+1900, time.Month(rawhdr.UpdateMonth), int(rawhdr.UpdateDay), 0, 0, 0, 0, time.Local),
		RecordCount: int(rawhdr.RecordCount),
		RecordSize:  int(rawhdr.LengthRecordBytes),
		FieldCount:  rawFieldCount,
	}

	if !isSupportedVersion(db.Header.Version) {
		return ErrorNotSupportedVersion(db.Header.Version)
	}

	return nil
}

// Read single terminator character (0x0d) after all field descriptors headers
func (db *DBaseFile) readTerminator() (err error) {
	err = db.checkOffset(db.offsets.fieldEnd, `starting to read terminator character`)
	if err != nil {
		return err
	}

	terminator := make([]byte, 1)
	_, err = db.r.Read(terminator)
	if err != nil {
		return err
	}

	if terminator[0] != 0x0d {
		return fmt.Errorf(`invalid terminator: %#v`, terminator)
	}

	err = db.checkOffset(db.offsets.terminatorEnd, `after terminator character`)
	if err != nil {
		return err
	}

	return nil
}
