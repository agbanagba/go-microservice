package files

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpLocal(t *testing.T) (*Local, string, func()) {
	dir, err := ioutil.TempDir("", "files")
	if err != nil {
		t.Fatal(err)
	}

	l, err := NewLocal(dir, 1024)
	if err != nil {
		t.Fatal(err)
	}

	return l, dir, func() {
		// cleanup
		os.RemoveAll(dir)
	}
}

func TestSaveContentsOfReader(t *testing.T) {
	savePath := "/1/test.png"
	fileContent := "Hello, World"

	l, dir, cleanup := setUpLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.NoError(t, err)

	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)

	// check the contents of the file
	d, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(d))
}

func TestGetsContentAndWritesToWriter(t *testing.T) {
	savePath := "/1/test.png"
	fileContent := "Hello, World"

	l, _, cleanup := setUpLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.NoError(t, err)

	// Read the file back
	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()

	d, err := ioutil.ReadAll(r)
	assert.Equal(t, fileContent, string(d))

}
