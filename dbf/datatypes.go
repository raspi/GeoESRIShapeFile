package dbf

import "fmt"

// Data types in FieldDescriptors
type DataType uint8

const (
	Character     DataType = 'C' // Text, String
	DateData      DataType = 'D' // Date, Time
	FloatingPoint DataType = 'F' // Float
	Numerical     DataType = 'N' // Decimal
	Logical       DataType = 'L' // Boolean
	MemoData      DataType = 'M' // Text?
)

func (dt DataType) String() string {
	switch dt {
	case Character:
		return "Character"
	case DateData:
		return "Date"
	case FloatingPoint:
		return "FloatingPoint"
	case Numerical:
		return "Numerical"
	case Logical:
		return "Logical"
	case MemoData:
		return "Memo"
	default:
		return fmt.Sprintf(`unknown: '%[1]c' %[1]d`, dt)
	}

}

func isSupportedDataType(d DataType) bool {
	switch d {
	case Character, DateData, FloatingPoint, Numerical, Logical, MemoData:
		return true
	default:
		return false
	}
}

type NotSupportedDataType struct {
	DataType DataType
}

func NewErrorNotSupportedDataType(d DataType) error {
	return NotSupportedDataType{DataType: d}.Error()
}

func (n NotSupportedDataType) Error() error {
	return fmt.Errorf(`not supported data type: %v`, n.DataType)
}
