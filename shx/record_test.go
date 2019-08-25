package shx

import (
	"encoding/binary"
	"testing"
)

func TestIndexRecordFileSize(t *testing.T) {
	actual := binary.Size(ShapeIndexRecord{})
	if actual != 8 {
		t.Fatalf(`record size was %v, should be 8`, actual)
	}
}
