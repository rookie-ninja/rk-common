// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

const (
	// OsR OS read
	OsR = 04
	// OsW OS write
	OsW = 02
	// OsX OS execute
	OsX = 01
	// OsUerShift OS user shift
	OsUerShift = 6
	// OsGroupShift OS group shift
	OsGroupShift = 3
	// OsOthShift OS OTH shift
	OsOthShift = 0
	// OsUserR OS user read
	OsUserR = OsR << OsUerShift
	// OsUserW OS user write
	OsUserW = OsW << OsUerShift
	// OsUserW OS user execute
	OsUserX = OsX << OsUerShift
	// OsUserRW OS user read and write
	OsUserRW = OsUserR | OsUserW
	// OsUserW OS user read, write and execute
	OsUserRWX = OsUserRW | OsUserX
	// OsGroupR OS group read
	OsGroupR = OsR << OsGroupShift
	// OsOthR OS OTH read
	OsOthR = OsR << OsOthShift
	// OsAllR OS all read
	OsAllR = OsUserR | OsGroupR | OsOthR
	// DefaultFileMode default file mode for RK
	DefaultFileMode = os.ModeDir | (OsUserRWX | OsAllR)
)

func NewMemFsToLocalFsCopier(memFs billy.Filesystem) *memFsToLocalFsCopier {
	return &memFsToLocalFsCopier{MemFs: memFs}
}

// memFsToLocalFsCopier Copy all files and directories from memory fs to local fs
type memFsToLocalFsCopier struct {
	MemFs billy.Filesystem
}

// CopyDir Copy directory from memory fs to local fs
func (c *memFsToLocalFsCopier) CopyDir(srcPath string, dstPath string) error {
	var err error
	var fds []os.FileInfo

	// 1: Create all directories at local fs destination at once
	if err = os.MkdirAll(dstPath, DefaultFileMode); err != nil {
		return err
	}

	// 2: Read directory from memory fs
	if fds, err = c.MemFs.ReadDir(srcPath); err != nil {
		return err
	}

	// 3: Iterate all files and directories and copy recursively
	for _, fd := range fds {
		srcFilePath := path.Join(srcPath, fd.Name())
		dstFilePath := path.Join(dstPath, fd.Name())

		if fd.IsDir() {
			// If it is a directory, then call CopyDir
			if err = c.CopyDir(srcFilePath, dstFilePath); err != nil {
				return err
			}
		} else {
			// If it is a file, then call CopyFile
			if err = c.CopyFile(srcFilePath, dstFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile Copy file from memory fs to local fs
func (c *memFsToLocalFsCopier) CopyFile(srcPath, dstPath string) error {
	var err error
	var srcFd billy.File
	var dstFd *os.File

	// 1: Open file from memory fs
	if srcFd, err = c.MemFs.Open(srcPath); err != nil {
		return err
	}
	defer srcFd.Close()

	// 2: Create file at local destination path
	if dstFd, err = os.Create(dstPath); err != nil {
		return err
	}
	defer dstFd.Close()

	// 3: Copy file content
	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}

	// 4: Grant file permission
	return os.Chmod(dstPath, DefaultFileMode)
}

// NewPkgerFsToLocalFsCopier Copy all files and directories from pkger to local fs
func NewPkgerFsToLocalFsCopier(pkgerIdentity string) *pkgerFsToLocalFsCopier {
	return &pkgerFsToLocalFsCopier{
		pkgerIdentity: pkgerIdentity,
	}
}

// pkgerFsToLocalFsCopier Copy all files and directories from memory fs to local fs
type pkgerFsToLocalFsCopier struct {
	pkgerIdentity string
}

// CopyDir Copy directory from memory fs to local fs
func (c *pkgerFsToLocalFsCopier) CopyDir(srcPath string, dstPath string) error {
	if !strings.HasPrefix(srcPath, "/") {
		srcPath = path.Join("/", srcPath)
	}

	var err error

	// 1: Create all directories at local fs destination at once
	if err = os.MkdirAll(dstPath, DefaultFileMode); err != nil {
		return err
	}

	// 2: Walk all files and create files and directories
	rootPath := fmt.Sprintf("%s:%s", c.pkgerIdentity, srcPath)
	if err = pkger.Walk(rootPath, func(currPath string, info fs.FileInfo, err error) error {
		// Get file path without identity
		subPath := strings.TrimPrefix(strings.TrimPrefix(currPath, rootPath), "/")

		// Root path, continue
		if len(subPath) < 1 {
			return nil
		}

		if info.IsDir() {
			// Create directory at destination
			if err = os.MkdirAll(path.Join(dstPath, subPath), DefaultFileMode); err != nil {
				return err
			}
		} else {
			// Copy file
			return c.CopyFile(currPath, path.Join(dstPath, subPath))
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// CopyFile Copy file from memory fs to local fs
func (c *pkgerFsToLocalFsCopier) CopyFile(srcPath, dstPath string) error {
	var err error
	var srcFd pkging.File
	var dstFd *os.File

	if !strings.Contains(srcPath, c.pkgerIdentity) {
		srcPath = fmt.Sprintf("%s:%s", c.pkgerIdentity, srcPath)
	}

	// 1: Open file from memory fs
	if srcFd, err = pkger.Open(srcPath); err != nil {
		return err
	}
	defer srcFd.Close()

	// 2: Create file at local destination path
	if dstFd, err = os.Create(dstPath); err != nil {
		return err
	}
	defer dstFd.Close()

	// 3: Copy file content
	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}

	// 4: Grant file permission
	return os.Chmod(dstPath, DefaultFileMode)
}
