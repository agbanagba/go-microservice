package files

import (
	"io"
)

// Storage defines
type Storage interface {
	Save(path string, file io.Reader) error
}
