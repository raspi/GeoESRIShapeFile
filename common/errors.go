package common

import (
	"errors"
	"fmt"
)

var ErrorNotInitialized = errors.New(`not initialized`)

type InvalidFileCode struct {
	Code uint32
}

func (e *InvalidFileCode) Error() string {
	return fmt.Sprintf(`invalid file code: %d`, e.Code)
}

type InvalidHeaderUnused struct {
	Index int
	Value uint32
}

func (e *InvalidHeaderUnused) Error() string {
	return fmt.Sprintf(`invalid unused %v at index: %d`, e.Value, e.Index)
}

type InvalidHeaderLength struct {
	Value uint32
}

func (e *InvalidHeaderLength) Error() string {
	return fmt.Sprintf(`invalid length: %d`, e.Value)
}

type InvalidHeaderVersion struct {
	Version uint32
}

func (e *InvalidHeaderVersion) Error() string {
	return fmt.Sprintf(`invalid version: %d`, e.Version)
}

type ErrInvalidShapeType struct {
	ShapeType ShapeType
}

func (e *ErrInvalidShapeType) Error() string {
	return fmt.Sprintf(`invalid shape type: %[1]d`, e.ShapeType)
}
