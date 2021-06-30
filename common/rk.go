// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
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
	RkMetaFilePath = ".rk/rk.yaml"
	RkDepFilePath = ".rk/dep/go.mod"
	RkUtHtmlFilePath = ".rk/ut/cov.html"
	RkUtOutFilepath = ".rk/ut/cov.out"
	RkLicenseFilePath = ".rk/LICENSE"
	RkReadmeFilePath = ".rk/README.md"
)

var (
	GetPackageNameArgs = []string{
		"rev-parse",
		"--show-toplevel",
	}

	GetCurrentTagArgs = []string{
		"tag",
		"--points-at",
		"HEAD",
	}

	GetRemoteUrlArgs = []string{
		"config",
		"--get",
		"remote.origin.url",
	}

	GetBranchArgs = []string{
		"rev-parse",
		"--abbrev-ref",
		"HEAD",
	}

	GetLatestCommitArgs = []string{
		"log",
		"-n1",
		`--pretty=format:{%n  "id": "%H",%n  "idAbbr": "%h",%n  "sub": "%s",%n  "date": "%cD",%n  "committer": {%n    "name": "%cN",%n    "email": "%cE"%n  }%n}`,
	}
)

type RkMeta struct {
	Name             string `json:"name" yaml:"name"`
	Version          string `json:"version" yaml:"version"`
	Git              *Git   `json:"git" yaml:"git"`
}

type Git struct {
	Url     string `yaml:"url" json:"url"`
	Branch  string `yaml:"branch" json:"branch"`
	Tag     string `yaml:"tag" json:"tag"`
	Commit  *Commit `yaml:"commit" json:"commit"`
}

type Commit struct {
	Id        string `yaml:"id" json:"id"`
	Date      string `yaml:"date" json:"date"`
	IdAbbr    string `yaml:"idAbbr" json:"idAbbr"`
	Sub       string `yaml:"sub" json:"sub"`
	Committer *Committer `yaml:"committer" json:"committer"`
}

type Committer struct {
	Name  string `yaml:"name" json:"name"`
	Email string `yaml:"email" json:"email"`
}

func GetPackageNameFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetPackageNameArgs...).CombinedOutput(); err != nil {
		return "", err
	} else {
		return path.Base(strings.TrimSuffix(string(v), "\n")), nil
	}
}

func GetCurrentTagFromGitLocal() (string, error) {
	if v, err := exec.Command("git", GetCurrentTagArgs...).CombinedOutput(); err != nil {
		return "", err
	} else {
		return strings.TrimSuffix(string(v), "\n"), nil
	}
}

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

type GoEnv struct {
	Ar          string `json:"AR"`
	Cc          string `json:"CC"`
	CgoCflags   string `json:"CGO_CFLAGS"`
	CgoCppflags string `json:"CGO_CPPFLAGS"`
	CgoCxxflags string `json:"CGO_CXXFLAGS"`
	CgoEnabled  string `json:"CGO_ENABLED"`
	CgoFflags   string `json:"CGO_FFLAGS"`
	CgoLdflags  string `json:"CGO_LDFLAGS"`
	Cxx         string `json:"CXX"`
	GccGo       string `json:"GCCGO"`
	Go111Module string `json:"GO111MODULE"`
	GoArch      string `json:"GOARCH"`
	GoBin       string `json:"GOBIN"`
	Gocache     string `json:"GOCACHE"`
	GoEnv       string `json:"GOENV"`
	GoExe       string `json:"GOEXE"`
	GoFlags     string `json:"GOFLAGS"`
	GoGccFlags  string `json:"GOGCCFLAGS"`
	GoHostArch  string `json:"GOHOSTARCH"`
	GoHostOs    string `json:"GOHOSTOS"`
	GoInsecure  string `json:"GOINSECURE"`
	GoMod       string `json:"GOMOD"`
	GoModCache  string `json:"GOMODCACHE"`
	GoNoProxy   string `json:"GONOPROXY"`
	GoNoSumDb   string `json:"GONOSUMDB"`
	GoOs        string `json:"GOOS"`
	GoPath      string `json:"GOPATH"`
	GoPrivate   string `json:"GOPRIVATE"`
	GoProxy     string `json:"GOPROXY"`
	GoRoot      string `json:"GOROOT"`
	GoSumDb     string `json:"GOSUMDB"`
	GoTmpDir    string `json:"GOTMPDIR"`
	GoToolDir   string `json:"GOTOOLDIR"`
	GoVcs       string `json:"GOVCS"`
	GoVersion   string `json:"GOVERSION"`
	PkgConfig   string `json:"PKG_CONFIG"`
}

func GetGoWd() string {
	return path.Dir(GetGoEnv().GoMod)
}

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

func GetGoEnv() *GoEnv {
	res := &GoEnv{}

	if v, err := exec.Command("go", "env", "-json").CombinedOutput(); err != nil {
		return res
	} else {
		json.Unmarshal(v, res)
	}

	return res
}

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