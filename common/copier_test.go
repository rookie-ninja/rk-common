// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestMemFsToLocalFsCopier_CopyDir(t *testing.T) {
	// make file as bellow in mem fs
	// .
	// ├── a.txt
	// ├── b
	// │   └── c.txt
	fs := memfs.New()
	fs.MkdirAll("b", DefaultFileMode)
	fs.Create("a.txt")
	fs.Create("b/c.txt")

	// Copy to testing temp dir
	tmpDir := t.TempDir()

	c := NewMemFsToLocalFsCopier(fs)
	c.CopyDir(".", tmpDir)

	assert.True(t, FileExists(path.Join(tmpDir, "a.txt")))
	assert.True(t, FileExists(path.Join(tmpDir, "b/c.txt")))
}

func TestMemFsToLocalFsCopier_CopyFile(t *testing.T) {
	// make file as bellow in mem fs
	// .
	// ├── a.txt
	fs := memfs.New()
	fs.MkdirAll("b", DefaultFileMode)
	fs.Create("a.txt")

	// Copy to testing temp dir
	tmpDir := t.TempDir()

	c := NewMemFsToLocalFsCopier(fs)
	c.CopyDir(".", tmpDir)

	assert.True(t, FileExists(path.Join(tmpDir, "a.txt")))
}

func TestPkgerFsToLocalFsCopier_CopyDir(t *testing.T) {
	// Copy to testing temp dir
	tmpDir := t.TempDir()

	c := NewPkgerFsToLocalFsCopier("github.com/rookie-ninja/rk-common")
	err := c.CopyDir("/common/testdata", tmpDir)
	assert.Nil(t, err)

	assert.True(t, FileExists(path.Join(tmpDir, "a.txt")))
	assert.True(t, FileExists(path.Join(tmpDir, "b/c.txt")))
}

func TestPkgerFsToLocalFsCopier_CopyFile(t *testing.T) {
	// Copy to testing temp dir
	tmpDir := t.TempDir()

	c := NewPkgerFsToLocalFsCopier("github.com/rookie-ninja/rk-common")
	err := c.CopyFile("/common/testdata/a.txt", path.Join(tmpDir, "a.txt"))
	assert.Nil(t, err)

	assert.True(t, FileExists(path.Join(tmpDir, "a.txt")))
}
