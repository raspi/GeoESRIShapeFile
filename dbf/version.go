package dbf

import (
	"fmt"
)

type NotSupportedVersion struct {
	Ver Version
}

func ErrorNotSupportedVersion(v Version) error {
	return NotSupportedVersion{Ver: v}.Error()
}

func (n NotSupportedVersion) Error() error {
	return fmt.Errorf(`not supported version: %v`, n.Ver)
}

type Version uint8

const (
	VerFoxBase1                           Version = 0x02
	VerdBASEIII                           Version = 0x03
	VerdBASEIIIwithMemo                   Version = 0x83
	VerVisualFoxPro                       Version = 0x30
	VerVisualFoxProWithAutoIncrement      Version = 0x31
	VerVisualFoxProWithVarcharOrVarbinary Version = 0x32
	VerdBASEIVSQLTableNoMemo              Version = 0x43
	VerdBASEIVSQLSystemNoMemo             Version = 0x63
	VerdBASEIVWithMemo                    Version = 0x8b
	VerdBASEIVSQLTableWithMemo            Version = 0xcb
	VerFoxPro2                            Version = 0xfb
	VerFoxPro2WithMemo                    Version = 0xf5
)

func (v Version) String() string {
	switch v {
	case VerFoxBase1:
		return "FoxBase 1.0"
	case VerdBASEIII:
		return "FoxBase 2.x / dBASE III"
	case VerdBASEIIIwithMemo:
		return "FoxBase 2.x / dBASE III with memo file"
	case VerVisualFoxPro:
		return "Visual FoxPro"
	case VerVisualFoxProWithAutoIncrement:
		return "Visual FoxPro with auto increment"
	case VerVisualFoxProWithVarcharOrVarbinary:
		return "Visual FoxPro with varchar/varbinary"
	case VerdBASEIVSQLTableNoMemo:
		return "dBASE IV SQL Table, no memo file"
	case VerdBASEIVSQLSystemNoMemo:
		return "dBASE IV SQL System, no memo file"
	case VerdBASEIVWithMemo:
		return "dBASE IV with memo file"
	case VerdBASEIVSQLTableWithMemo:
		return "dBASE IV SQL Table with memo file"
	case VerFoxPro2:
		return "FoxPro 2"
	case VerFoxPro2WithMemo:
		return "FoxPro 2 with memo file"
	default:
		return fmt.Sprintf(`Unknown version: %d`, v)
	}
}

// Supported dBase versions
func isSupportedVersion(v Version) bool {
	switch v {
	case VerdBASEIII:
		return true
	default:
		return false
	}
}
