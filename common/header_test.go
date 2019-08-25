package common

import (
	"encoding/binary"
	"io"
	"testing"
)

func TestFirstHeaderBE(t *testing.T) {
	actual := binary.Size(ShapeFileHeader1{})
	if actual != 28 {
		t.Fatalf(`header size was %v, should be 28`, actual)
	}
	t.Logf(`actual size was %v`, actual)
}

func TestSecondHeaderLE(t *testing.T) {
	actual := binary.Size(ShapeFileHeader2{})
	if actual != 72 {
		t.Fatalf(`header size was %v, should be 72`, actual)
	}
	t.Logf(`actual size was %v`, actual)
}

func TestFirstHeaderBEWithoutData(t *testing.T) {
	data := []byte{}

	r, err := NewReadSeekCloser(data)
	if err != nil {
		t.Fatal(err)
	}

	err = readFirstHeader(r)
	if err != io.EOF {
		t.Fatal(err)
	}
}

func TestFirstHeaderBEWithAllZero(t *testing.T) {
	data := []byte{
		0x0, 0x0, 0x0, 0x0, // File code
		0x0, 0x0, 0x0, 0x0, // unused 0
		0x0, 0x0, 0x0, 0x0, // unused 1
		0x0, 0x0, 0x0, 0x0, // unused 2
		0x0, 0x0, 0x0, 0x0, // unused 3
		0x0, 0x0, 0x0, 0x0, // unused 4
		0x0, 0x0, 0x0, 0x0, // length
	}

	r, err := NewReadSeekCloser(data)
	if err != nil {
		t.Fatal(err)
	}

	expectederr := InvalidFileCode{Code: 0}
	err = readFirstHeader(r)

	convertederr, ok := err.(*InvalidFileCode)
	if !ok {
		t.Fatalf(`error was %v instead of %v and failed type conversion`, err, expectederr)
	}

	if *convertederr != expectederr {
		t.Fatalf(`error %v instead of %v`, convertederr, expectederr)
	}

	t.Logf(`error was: %#v`, convertederr)
}

func TestFirstHeaderBEWithFileCode(t *testing.T) {
	data := []byte{
		0x67, 0xff, 0x80, 0x12, // File code
		0x0, 0x0, 0x0, 0x0, // unused 0
		0x0, 0x0, 0x0, 0x0, // unused 1
		0x0, 0x0, 0x0, 0x0, // unused 2
		0x0, 0x0, 0x0, 0x0, // unused 3
		0x0, 0x0, 0x0, 0x0, // unused 4
		0x0, 0x0, 0x0, 0x0, // length
	}

	r, err := NewReadSeekCloser(data)
	if err != nil {
		t.Fatal(err)
	}

	expectederr := InvalidFileCode{Code: 1744797714}
	err = readFirstHeader(r)

	convertederr, ok := err.(*InvalidFileCode)
	if !ok {
		t.Fatalf(`error was %v instead of %v and failed type conversion`, err, expectederr)
	}

	if *convertederr != expectederr {
		t.Fatalf(`error %v instead of %v`, convertederr, expectederr)
	}

	t.Logf(`error was: %#v`, convertederr)
}

func TestFirstHeaderBEWithCorrectFileCode(t *testing.T) {
	data := []byte{
		0x0, 0x0, 0x27, 0x0a, // File code
		0x8, 0xf, 0xa, 0x4, // unused 0
		0xff, 0xff, 0xff, 0xff, // unused 1
		0xff, 0xff, 0xff, 0xff, // unused 2
		0xff, 0xff, 0xff, 0xff, // unused 3
		0xff, 0xff, 0xff, 0xff, // unused 4
		0x0, 0x0, 0x0, 0x0, // length
	}

	r, err := NewReadSeekCloser(data)
	if err != nil {
		t.Fatal(err)
	}

	expectederr := InvalidHeaderUnused{Index: 0, Value: 135203332}
	err = readFirstHeader(r)

	convertederr, ok := err.(*InvalidHeaderUnused)
	if !ok {
		t.Fatalf(`error was %v instead of %v and failed type conversion`, err, expectederr)
	}

	if *convertederr != expectederr {
		t.Fatalf(`error %v instead of %v`, convertederr, expectederr)
	}

	t.Logf(`error was: %#v`, convertederr)
}

func TestFirstHeaderBEWithBrokenLength(t *testing.T) {
	data := []byte{
		0x0, 0x0, 0x27, 0x0a, // File code
		0x0, 0x0, 0x0, 0x0, // unused 0
		0x0, 0x0, 0x0, 0x0, // unused 1
		0x0, 0x0, 0x0, 0x0, // unused 2
		0x0, 0x0, 0x0, 0x0, // unused 3
		0x0, 0x0, 0x0, 0x0, // unused 4
		0x0, 0x0, 0x0, 0x0, // length
	}

	r, err := NewReadSeekCloser(data)
	if err != nil {
		t.Fatal(err)
	}

	expectederr := InvalidHeaderLength{Value: 0}
	err = readFirstHeader(r)

	convertederr, ok := err.(*InvalidHeaderLength)
	if !ok {
		t.Fatalf(`error was %v instead of %v and failed type conversion`, err, expectederr)
	}

	if *convertederr != expectederr {
		t.Fatalf(`error %v instead of %v`, convertederr, expectederr)
	}

	t.Logf(`error was: %#v`, convertederr)
}
