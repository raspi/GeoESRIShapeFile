package common

import (
	"github.com/spf13/afero"
	"io"
	"time"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
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
