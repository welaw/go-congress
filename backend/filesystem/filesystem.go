package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/welaw/go-congress/pkg/errs"
)

type Filesystem interface {
	OpenFile(string) ([]byte, error)
	WriteFile(string, []byte) error
}

type filesystem struct {
	path string
}

func NewFilesystem(path string) Filesystem {
	return &filesystem{
		path: path,
	}
}

func (fs *filesystem) OpenFile(path string) ([]byte, error) {
	// check if exists
	full := fmt.Sprintf("%s/%s", fs.path, path)
	_, err := os.Stat(full)
	if os.IsNotExist(err) {
		return nil, errs.ErrNotFound
	}

	// read file
	file, err := ioutil.ReadFile(full)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fs *filesystem) WriteFile(path string, contents []byte) error {
	full := fmt.Sprintf("%s/%s", fs.path, path)

	if _, err := os.Stat(fs.path); os.IsNotExist(err) {
		os.MkdirAll(fs.path, 755)
	}

	err := ioutil.WriteFile(full, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
