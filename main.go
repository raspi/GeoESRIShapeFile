package geoesrishapefile

import (
	"github.com/raspi/GeoESRIShapeFile/dbf"
	"github.com/raspi/GeoESRIShapeFile/shp"
	"github.com/raspi/GeoESRIShapeFile/shx"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ShapeFiles struct {
	Fshp  shp.ShapeFile
	Fshx  []shx.ShapeOffsetIndex
	Fdbf  dbf.DBaseFile
	debug bool
}

func New(fpath string, parseFieldNames []string, parseFieldNamesOperation dbf.Operation, defaultConverter dbf.ConverterFunction, converters map[string]dbf.ConverterFunction) (sf ShapeFiles, err error) {
	sf.debug = true

	fpath, err = filepath.Abs(fpath)
	if err != nil {
		return sf, err
	}

	_, err = os.Stat(fpath)
	if err != nil {
		return sf, err
	}

	origdir, origfname := filepath.Split(fpath)
	origext := filepath.Ext(origfname)
	origFnameNoExt := strings.TrimRight(origfname, origext)

	flist, err := ioutil.ReadDir(origdir)
	if err != nil {
		return sf, err
	}

	for _, f := range flist {
		if f.IsDir() {
			continue
		}

		dir, fname, ext := fnamesplit(filepath.Join(origdir, f.Name()))

		if fname != origFnameNoExt {
			continue
		}

		ofile := filepath.Join(dir, fname) + "." + ext

		switch strings.ToLower(ext) {
		case `dbf`:
			sf.Fdbf, err = dbf.New(ofile, parseFieldNames, parseFieldNamesOperation, defaultConverter, converters)
			if err != nil {
				return sf, err
			}
		case `shp`:
			sf.Fshp, err = shp.New(ofile)
			if err != nil {
				return sf, err
			}
		case `shx`:
			sf.Fshx, _, err = shx.New(ofile)
			if err != nil {
				return sf, err
			}

		default:
			continue
		}
	}

	return sf, nil
}

func fnamesplit(fpath string) (dir, fname, ext string) {
	dir, fname = filepath.Split(fpath)
	ext = filepath.Ext(fname)
	fnameNoExt := strings.TrimRight(fname, ext)
	ext = strings.TrimLeft(ext, `.`)
	return dir, fnameNoExt, ext
}
