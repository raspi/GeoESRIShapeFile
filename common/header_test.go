package common

import (
	"encoding/binary"
	"testing"
)

func TestFirstHeaderBE(t *testing.T) {
	actual := binary.Size(ShapeFileHeader1{})
	if actual != 28 {
		t.Fatalf(`header size was %v, should be 28`, actual)
	}
}

func TestSecondHeaderLE(t *testing.T) {
	actual := binary.Size(ShapeFileHeader2{})
	if actual != 72 {
		t.Fatalf(`header size was %v, should be 72`, actual)
	}
}
