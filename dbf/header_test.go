package dbf

import (
	"encoding/binary"
	"testing"
)

func TestRawHeaderSize(t *testing.T) {
	actual := binary.Size(rawHeader{})
	if actual != 32 {
		t.Fatalf(`header size was %v, should be 32`, actual)
	}
}

func TestRawFieldHeaderSize(t *testing.T) {
	actual := binary.Size(rawFieldDescriptor{})
	if actual != 32 {
		t.Fatalf(`field header size was %v, should be 32`, actual)
	}
}
