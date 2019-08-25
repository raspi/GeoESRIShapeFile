package shp

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/xerrors"
	"io"
)

type ShapeTypeI interface {
	Validate() error
	read(r io.ReadSeeker) (ShapeTypeI, error)
}

type Box struct {
	MinX, MinY, MaxX, MaxY float64
}

func (b Box) String() string {
	return fmt.Sprintf(`%f, %f x %f, %f`, b.MinX, b.MaxX, b.MinY, b.MaxY)
}

type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf(`%v x %v`, p.X, p.Y)
}

type PolyLineZ struct {
	Box       Box
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Points    []Point
	ZRange    [2]float64
	ZArray    []float64
	MRange    [2]float64
	MArray    []float64
}

func (p PolyLineZ) String() string {
	return fmt.Sprintf(`%v parts %v points Box(%v)`, p.NumParts, p.NumPoints, p.Box)
}

func (z PolyLineZ) Validate() error {
	if len(z.Points) != int(z.NumPoints) {
		return fmt.Errorf(`numpoints mismatch`)
	}

	if len(z.Parts) != int(z.NumParts) {
		return fmt.Errorf(`numparts mismatch`)
	}

	return nil
}

/*
	Position     Field      Value     Type    Number    Order
	Byte 0       Shape Type 13        Integer 1         Little
	Byte 4       Box        Box       Double  4         Little
	Byte 36      NumParts   NumParts  Integer 1         Little
	Byte 40      NumPoints  NumPoints Integer 1         Little
	Byte 44      Parts      Parts     Integer NumParts  Little
	Byte X       Points     Points    Point   NumPoints Little
	Byte Y       Zmin       Zmin      Double  1         Little
	Byte Y + 8   Zmax       Zmax      Double  1         Little
	Byte Y + 16  Zarray     Zarray    Double  NumPoints Little
	Byte Z*      Mmin       Mmin      Double  1         Little
	Byte Z + 8*  Mmax       Mmax      Double  1         Little
	Byte Z + 16* Marray     Marray    Double  NumPoints Little

	Note:  X = 44 + (4 * NumParts), Y = X + (16 * NumPoints), Z = Y + 16 + (8 * NumPoints)*  optional
*/
func (z PolyLineZ) read(r io.ReadSeeker) (retrec ShapeTypeI, err error) {
	var rawrec struct {
		Box       Box
		NumParts  uint32
		NumPoints uint32
	}
	err = binary.Read(r, binary.LittleEndian, &rawrec)
	if err != nil {
		return nil, xerrors.Errorf(`couldn't read raw poly line header: %w`, err)
	}

	z.Box = rawrec.Box
	z.NumParts = rawrec.NumParts
	z.NumPoints = rawrec.NumPoints

	z.Parts = make([]uint32, z.NumParts)
	err = binary.Read(r, binary.LittleEndian, &z.Parts)
	if err != nil {
		return nil, xerrors.Errorf(`parts: %w`, err)
	}

	z.Points = make([]Point, z.NumPoints)
	err = binary.Read(r, binary.LittleEndian, &z.Points)
	if err != nil {
		return nil, xerrors.Errorf(`points: %w`, err)
	}

	err = binary.Read(r, binary.LittleEndian, &z.ZRange)
	if err != nil {
		return nil, xerrors.Errorf(`zrange: %w`, err)
	}

	z.ZArray = make([]float64, z.NumPoints)
	err = binary.Read(r, binary.LittleEndian, &z.ZArray)
	if err != nil {
		return nil, xerrors.Errorf(`Z-Array: %w`, err)
	}

	err = binary.Read(r, binary.LittleEndian, &z.MRange)
	if err != nil {
		return nil, xerrors.Errorf(`M Range: %w`, err)
	}

	z.MArray = make([]float64, z.NumPoints)
	err = binary.Read(r, binary.LittleEndian, &z.MArray)
	if err != nil {
		return nil, xerrors.Errorf(`M-Array: %w`, err)
	}

	return z, nil
}
