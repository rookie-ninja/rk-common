// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func TestGetBootConfigPath_WithFlags(t *testing.T) {
	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(""), 0777))

	GlobalFlags.Set("rkboot", filePath)
	defer GlobalFlags.Set("rkboot", "")

	assert.Equal(t, filePath, GetBootConfigPath(""))
}

func TestGetBootConfigPath_WithFlagsAndRelativePath(t *testing.T) {
	// Use flags.go files located on current working directory.
	filePath := "flags.go"
	GlobalFlags.Set("rkboot", filePath)
	defer GlobalFlags.Set("rkboot", "")
	assert.Contains(t, GetBootConfigPath(""), filePath)
}

func TestGetBootConfigPath_WithNonExistFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()

	filePath := "non-exist.yaml"
	GlobalFlags.Set("rkboot", filePath)
	defer GlobalFlags.Set("rkboot", "")

	GetBootConfigPath("")
}

func TestGetBootConfigPath_WithoutOverride(t *testing.T) {
	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(""), 0777))
	assert.Equal(t, filePath, GetBootConfigPath(filePath))
}

func TestGetBootConfigOverrides_WithoutOverrides(t *testing.T) {
	assert.Empty(t, GetBootConfigOverrides())
}

func TestGetBootConfigOverrides_ExpectPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()

	GlobalFlags.Set("rkset", "should panic")
	defer GlobalFlags.Set("rkset", "")

	GetBootConfigOverrides()
}

func TestGetBootConfigOverrides_HappyCase(t *testing.T) {
	GlobalFlags.Set("rkset", "key=value,slice[0]=value")
	defer GlobalFlags.Set("rkset", "")

	res := GetBootConfigOverrides()
	assert.Equal(t, "value", res["key"])
	assert.Equal(t, "value", res["slice"].([]interface{})[0])
}

func TestGetBootConfigOriginal_WithNonExistFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()

	GetBootConfigOriginal("non-exist")
}

func TestGetBootConfigOriginal_WithInvalidFileType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()

	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte("should panic"), 0777))

	GetBootConfigOriginal(filePath)
}

func TestGetBootConfigOriginal_HappyCase(t *testing.T) {
	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(`
---
key: value
`), 0777))

	res := GetBootConfigOriginal(filePath)
	assert.Equal(t, "value", res["key"])
}

func TestUnmarshalBootConfig_HappyCase(t *testing.T) {
	// Set flags in order to override
	GlobalFlags.Set("rkset", "key=value2")
	defer GlobalFlags.Set("rkset", "")

	// Write original file to local file system
	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(`
---
key: value
`), 0777))

	type MyStruct struct {
		Key string `yaml:"key"`
	}
	config := &MyStruct{}
	UnmarshalBootConfig(filePath, config)

	// Value should be overridden
	assert.Equal(t, "value2", config.Key)
}

func TestUnmarshalBootConfig_WithInvalidConfigType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// expect panic to be called with non nil error
			assert.True(t, true)
		} else {
			// this should never be called in case of a bug
			assert.True(t, false)
		}
	}()

	// Write original file to local file system
	filePath := path.Join(t.TempDir(), "ut-temp.yaml")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte(`
---
key: value
`), 0777))

	UnmarshalBootConfig(filePath, "should panic")
}
