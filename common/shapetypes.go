package common

import "fmt"

type ShapeType int32

const (
	NULL        ShapeType = 0
	POINT       ShapeType = 1
	POLYLINE    ShapeType = 3
	POLYGON     ShapeType = 5
	MULTIPOINT  ShapeType = 8
	POINTZ      ShapeType = 11
	POLYLINEZ   ShapeType = 13
	POLYGONZ    ShapeType = 15
	MULTIPOINTZ ShapeType = 18
	POINTM      ShapeType = 21
	POLYLINEM   ShapeType = 23
	POLYGONM    ShapeType = 25
	MULTIPOINTM ShapeType = 28
	MULTIPATCH  ShapeType = 31
)

func (st ShapeType) String() string {
	switch st {
	case NULL:
		return "Null"
	case POINT:
		return "Point"
	case POLYLINE:
		return "PolyLine"
	case POLYGON:
		return "Polygon"
	case MULTIPOINT:
		return "MultiPoint"
	case POINTZ:
		return "PointZ"
	case POLYLINEZ:
		return "PolyLineZ"
	case POLYGONZ:
		return "PolygonZ"
	case MULTIPOINTZ:
		return "MultiPointZ"
	case POINTM:
		return "PointM"
	case POLYLINEM:
		return "PolyLineM"
	case POLYGONM:
		return "PolygonM"
	case MULTIPOINTM:
		return "MultiPointM"
	case MULTIPATCH:
		return "MultiPatch"
	default:
		return fmt.Sprintf("Unsupported shape type: %d", st)
	}
}

func IsSupportedShapeType(st ShapeType) bool {
	switch st {
	case NULL, POINT, POLYLINE, POLYGON, MULTIPOINT, POINTZ, POLYLINEZ, POLYGONZ, MULTIPOINTZ, POINTM, POLYLINEM, POLYGONM, MULTIPOINTM, MULTIPATCH:
		return true
	default:
		return false
	}
}
