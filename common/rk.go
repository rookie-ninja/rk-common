// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

import (
	"encoding/json"
	"github.com/rogpeppe/go-internal/modfile"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

const (
	// RkMetaFilePath used in rk cli
	RkMetaFilePath = ".rk/rk.yaml"
	// RkDepFilePath used in rk cli
	RkDepFilePath = ".rk/dep/go.mod"
	// RkUtHtmlFilePath used in rk cli
	RkUtHtmlFilePath = ".rk/ut/cov.html"
	// RkUtOutFilepath used in rk cli
	RkUtOutFilepath = ".rk/ut/cov.out"
	// RkLicenseFilePath used in rk cli
	RkLicenseFilePath = ".rk/LICENSE"
	// RkReadmeFilePath used in rk cli
	RkReadmeFilePath = ".rk/README.md"
)

var (
	// GetPackageNameArgs used in rk cli for reading local git meta info
	GetPackageNameArgs = []string{
		"rev-parse",
		"--show-toplevel",
	}
	// GetCurrentTagArgs used in rk cli for reading local git meta info
	GetCurrentTagArgs = []string{
		"tag",
		"--points-at",
		"HEAD",
	}
	// GetRemoteUrlArgs used in rk cli for reading local git meta info
	GetRemoteUrlArgs = []string{
		"config",
		"--get",
		"remote.origin.url",
	}
	// GetBranchArgs used in rk cli for reading local git meta info
	GetBranchArgs = []string{
		"rev-parse",
		"--abbrev-ref",
		"HEAD",
	}
	// GetLatestCommitArgs used in rk cli for reading local git meta info
	GetLatestCommitArgs = []string{
		"log",
		"-n1",
		`--pretty=format:{%n  "id": "%H",%n  "idAbbr": "%h",%n  "sub": "%s",%n  "date": "%cD",%n  "committer": {%n    "name": "%cN",%n    "email": "%cE"%n  }%n}`,
	}
)

// RkMeta would be extracted by rk cli
type RkMeta struct {
	// Name of application
	Name string `json:"name" yaml:"name"`
	// Version of application
	Version string `json:"version" yaml:"version"`
	// Git meta info
	Git *Git `json:"git" yaml:"git"`
}

// Git metadata info on local machine
type Git struct {
	// Url of git repo
	Url string `yaml:"url" json:"url"`
	// Branch of git repo
	Branch string `yaml:"branch" json:"branch"`
	// Tag of git repo
	Tag string `yaml:"tag" json:"tag"`
	// Commit info of git repo
	Commit *Commit `yaml:"commit" json:"commit"`
}

// Commit of git from local machine
type Commit struct {
	// Id of current commit
	Id string `yaml:"id" json:"id"`
	// Date of current commit
	Date string `yaml:"date" json:"date"`
	// IdAbbr is abbreviation of id of current commit
	IdAbbr string `yaml:"idAbbr" json:"idAbbr"`
	// Sub is subject of current commit
	Sub string `yaml:"sub" json:"sub"`
	// Committer of current commit
	Committer *Committer `yaml:"committer" json:"committer"`
}

// Committer info of current commit
type Committer struct {
	// Name of committer
	Name string `yaml:"name" json:"name"`
	// Email of committer
	Email string `yaml:"email" json:"email"`
}

// GetPackageNameFromGitLocal returns package name from git info on local machine
func GetPackageNameFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetPackageNameArgs...).CombinedOutput(); err != nil {
		return "", err
	} else {
		return path.Base(strings.TrimSuffix(string(v), "\n")), nil
	}
}

// GetCurrentTagFromGitLocal returns current tag info from git info on local machine
func GetCurrentTagFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetCurrentTagArgs...).CombinedOutput(); err != nil {
		return "", err
	} else {
		return strings.TrimSuffix(string(v), "\n"), nil
	}
}

// GetRemoteUrlFromGitLocal returns remote url from git info on local machine
func GetRemoteUrlFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetRemoteUrlArgs...).CombinedOutput(); err != nil {
		return "", err
	} else {
		rawUrl := strings.TrimSuffix(strings.TrimSuffix(string(v), "\n"), ".git")
		if strings.HasPrefix(rawUrl, "https") || strings.HasPrefix(rawUrl, "http") {
			return rawUrl, nil
		} else {
			tokens := strings.SplitN(rawUrl, "@", 2)
			if len(tokens) == 2 {
				return strings.TrimSuffix(tokens[1], "\n"), nil
			}
		}
	}

	return "", nil
}

// GetBranchFromGitLocal returns branch from git info on local machine
func GetBranchFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetBranchArgs...).CombinedOutput(); err != nil {
		return "", nil
	} else {
		branch := strings.TrimSuffix(string(v), "\n")
		if branch != "HEAD" {
			return branch, nil
		}
	}

	return "", nil
}

// GetLatestCommitFromGitLocal returns latest commit info from git info on local machine
func GetLatestCommitFromGitLocal() (*Commit, error) {
	res := &Commit{}

	if v, err := exec.Command("git", GetLatestCommitArgs...).CombinedOutput(); err != nil {
		return res, err
	} else {
		if err := json.Unmarshal(v, res); err != nil {
			return res, err
		}
	}

	return res, nil
}

// GoEnv defines golang environment variables
type GoEnv struct {
	Ar          string `json:"AR"`           // Ar from Go environment
	Cc          string `json:"CC"`           // Cc from Go environment
	CgoCflags   string `json:"CGO_CFLAGS"`   // CgoCflags from Go environment
	CgoCppflags string `json:"CGO_CPPFLAGS"` // CgoCppflags from Go environment
	CgoCxxflags string `json:"CGO_CXXFLAGS"` // CgoCxxflags from Go environment
	CgoEnabled  string `json:"CGO_ENABLED"`  // CgoEnabled from Go environment
	CgoFflags   string `json:"CGO_FFLAGS"`   // CgoFflags from Go environment
	CgoLdflags  string `json:"CGO_LDFLAGS"`  // CgoLdflags from Go environment
	Cxx         string `json:"CXX"`          // Cxx from Go environment
	GccGo       string `json:"GCCGO"`        // GccGo from Go environment
	Go111Module string `json:"GO111MODULE"`  // Go111Module from Go environment
	GoArch      string `json:"GOARCH"`       // GoArch from Go environment
	GoBin       string `json:"GOBIN"`        // GoBin from Go environment
	Gocache     string `json:"GOCACHE"`      // Gocache from Go environment
	GoEnv       string `json:"GOENV"`        // GoEnv from Go environment
	GoExe       string `json:"GOEXE"`        // GoExe from Go environment
	GoFlags     string `json:"GOFLAGS"`      // GoFlags from Go environment
	GoGccFlags  string `json:"GOGCCFLAGS"`   // GoGccFlags from Go environment
	GoHostArch  string `json:"GOHOSTARCH"`   // GoHostArch from Go environment
	GoHostOs    string `json:"GOHOSTOS"`     // GoHostOs from Go environment
	GoInsecure  string `json:"GOINSECURE"`   // GoInsecure from Go environment
	GoMod       string `json:"GOMOD"`        // GoMod from Go environment
	GoModCache  string `json:"GOMODCACHE"`   // GoModCache from Go environment
	GoNoProxy   string `json:"GONOPROXY"`    // GoNoProxy from Go environment
	GoNoSumDb   string `json:"GONOSUMDB"`    // GoNoSumDb from Go environment
	GoOs        string `json:"GOOS"`         // GoOs from Go environment
	GoPath      string `json:"GOPATH"`       // GoPath from Go environment
	GoPrivate   string `json:"GOPRIVATE"`    // GoPrivate from Go environment
	GoProxy     string `json:"GOPROXY"`      // GoProxy from Go environment
	GoRoot      string `json:"GOROOT"`       // GoRoot from Go environment
	GoSumDb     string `json:"GOSUMDB"`      // GoSumDb from Go environment
	GoTmpDir    string `json:"GOTMPDIR"`     // GoTmpDir from Go environment
	GoToolDir   string `json:"GOTOOLDIR"`    // GoToolDir from Go environment
	GoVcs       string `json:"GOVCS"`        // GoVcs from Go environment
	GoVersion   string `json:"GOVERSION"`    // GoVersion from Go environment
	PkgConfig   string `json:"PKG_CONFIG"`   // PkgConfig from Go environment
}

// GetGoWd return current working directory of golang
func GetGoWd() string {
	return path.Dir(GetGoEnv().GoMod)
}

// GetGoPkgName returns package name redden from go environment
func GetGoPkgName() string {
	res := ""
	modFilePath := GetGoEnv().GoMod

	if len(modFilePath) < 1 {
		return res
	}

	// read and extract application name
	if raw, err := ioutil.ReadFile(modFilePath); err != nil {
		return res
	} else {
		res = path.Base(modfile.ModulePath(raw))
	}

	return res

	return path.Base(GetGoEnv().GoMod)
}

// GetGoEnv returns GoEnv variable
func GetGoEnv() *GoEnv {
	res := &GoEnv{}

	if v, err := exec.Command("go", "env", "-json").CombinedOutput(); err != nil {
		return res
	} else {
		json.Unmarshal(v, res)
	}

	return res
}

// GetRkMetaFromCmd construct RkMeta from local environment
func GetRkMetaFromCmd() *RkMeta {
	meta := &RkMeta{}

	// Load package name from go.mod file
	meta.Name = GetGoPkgName()

	// Load current git info from local git
	meta.Git = &Git{}
	meta.Git.Branch, _ = GetBranchFromGitLocal()
	meta.Git.Tag, _ = GetCurrentTagFromGitLocal()
	meta.Git.Url, _ = GetRemoteUrlFromGitLocal()

	// Load latest commit from local git
	meta.Git.Commit, _ = GetLatestCommitFromGitLocal()

	// Load version
	if len(meta.Git.Tag) < 1 {
		// Tag is empty, let's try to construct version as <branch>-<commit IdAbbr>
		if len(meta.Git.Branch) < 1 {
			// Branch is empty, we will use commit IdAbbr instead
			meta.Version = meta.Git.Commit.IdAbbr
		} else {
			// <branch>-<commit IdAbbr>
			meta.Version = strings.Join([]string{meta.Git.Branch, meta.Git.Commit.IdAbbr}, "-")
		}
	} else {
		// Use tag as version
		meta.Version = meta.Git.Tag
	}

	return meta
}
