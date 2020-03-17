package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Local  is an implementation of the storage interface
type Local struct {
	maxFileSize int // max number of bytes
	basePath    string
}

// NewLocal creates a new Local filesystem with the base path.
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

// Save contents of the writer to the path
func (l *Local) Save(path string, contents io.Reader) error {
	fp := l.fullpath(path)

	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// if the file exists, delete it
	_, err = os.Stat(fp)
	if err == nil {
		err := os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to get file: %w", err)
	}

	// create the file at the path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil

}

// return absolute path
func (l *Local) fullpath(path string) string {
	return filepath.Join(l.basePath, path)
}

// Get the file at the given path and return a reader
func (l *Local) Get(path string) (*os.File, error) {
	fp := l.fullpath(path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file %w", err)
	}
	return f, nil
}
