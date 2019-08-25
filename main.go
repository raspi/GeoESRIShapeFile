package geoesrishapefile

import (
	"github.com/raspi/GeoESRIShapeFile/dbf"
	"github.com/raspi/GeoESRIShapeFile/shp"
	"github.com/raspi/GeoESRIShapeFile/shx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ShapeFiles struct {
	Fshp shp.ShapeFile
	Fshx shx.ShapeFileIndex // lookups
	Fdbf dbf.DBaseFile

	// logging
	debug struct {
		shp  bool
		shx  bool
		dbf  bool
		self bool
		all  bool
	}
}

func New(fpath string, parseFieldNames []string, parseFieldNamesOperation dbf.Operation, defaultConverter dbf.ConverterFunction, converters map[string]dbf.ConverterFunction) (sf ShapeFiles, err error) {
	sf.debug.all = true

	if sf.debug.all {
		sf.debug.self = true
		sf.debug.shp = true
		sf.debug.shx = true
		sf.debug.dbf = true
	}

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
		case `dbf`: // dBase Database
			err = sf.loadDbf(ofile, parseFieldNames, parseFieldNamesOperation, defaultConverter, converters)
			if err != nil {
				return sf, err
			}
		case `shp`: // ShapeFile
			err = sf.loadShp(ofile)
			if err != nil {
				return sf, err
			}
		case `shx`: // ShapeFile Index Offsets
			err = sf.loadShx(ofile)
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

func (sf *ShapeFiles) loadShp(fname string) (err error) {
	if sf.debug.shp {
		log.Printf(`loading .shp file %v`, fname)
	}

	sf.Fshp, err = shp.New(fname)
	if err != nil {
		return err
	}

	sf.Fshp.SetDebug(sf.debug.shp)

	err = sf.Fshp.Initialize()
	if err != nil {
		return err
	}

	return nil
}
func (sf *ShapeFiles) loadShx(fname string) (err error) {
	if sf.debug.shx {
		log.Printf(`loading .shx file %v`, fname)
	}

	sf.Fshx, err = shx.New(fname)
	if err != nil {
		return err
	}
	sf.Fshx.SetDebug(sf.debug.shx)

	err = sf.Fshx.Initialize()
	if err != nil {
		return err
	}
	return nil
}

func (sf *ShapeFiles) loadDbf(fname string, parseFieldNames []string, parseFieldNamesOperation dbf.Operation, defaultConverter dbf.ConverterFunction, converters map[string]dbf.ConverterFunction) (err error) {
	if sf.debug.dbf {
		log.Printf(`loading .dbf file %v`, fname)
	}

	sf.Fdbf, err = dbf.New(fname, parseFieldNames, parseFieldNamesOperation, defaultConverter, converters)
	if err != nil {
		return err
	}

	sf.Fdbf.SetDebug(sf.debug.dbf)

	err = sf.Fdbf.Initialize()
	if err != nil {
		return err
	}

	return nil
}
