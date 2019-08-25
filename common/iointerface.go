package common

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"time"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

func NewReadSeekCloser(b []byte) (rsc ReadSeekCloser, err error) {
	base := afero.NewMemMapFs()
	f, err := afero.TempFile(base, `/tmp`, `tmp_`)
	if err != nil {
		return f, err
	}

	wBytes, err := f.Write(b)
	if err != nil {
		return f, err
	}

	if wBytes != len(b) {
		return f, fmt.Errorf(`couldn't write`)
	}

	offset, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return f, err
	}

	if offset != 0 {
		return f, fmt.Errorf(`couldn't seek to start`)
	}

	return f, nil
}

func OpenFile(fpath string) (ReadSeekCloser, error) {
	base := afero.NewOsFs()
	layer := afero.NewMemMapFs()
	ufs := afero.NewCacheOnReadFs(base, layer, 5*time.Minute)

	f, err := ufs.Open(fpath)
	if err != nil {
		return nil, err
	}

	return f, nil
}
