package rk_common

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestShutdownWithError_WithNilError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()
	ShutdownWithError(nil)
}

func TestShutdownWithError_HappyCase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()
	ShutdownWithError(errors.New("error from unit test"))
}

func TestMustReadFile_WithEmptyFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()
	// expect panic
	MustReadFile("")
}

func TestMustReadFile_WithNonExistRelativeFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()
	// expect panic
	MustReadFile("non-exist-file-hopefully")
}

func TestMustReadFile_WithNonExistAbsoluteFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()
	// expect panic
	MustReadFile("/non-exist-file-hopefully")
}

func TestMustReadFile_HappyCase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	filePath := path.Join(t.TempDir(), "tmp-file")
	text := "some text"
	// write a file
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(text), os.ModePerm))
	assert.Equal(t, text, string(MustReadFile(filePath)))
}
