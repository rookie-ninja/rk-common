// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkcommon

import (
	"encoding/json"
	"os/exec"
	"path"
	"strings"
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
		`--pretty=format:{%n  "id": "%H",%n  "abbr": "%h",%n  "sub": "%s",%n  "date": "%cD",%n  "committer": {%n    "name": "%cN",%n    "email": "%cE"%n  }%n}`,
	}
)

type GitInfo struct {
	Package string `yaml:"package" json:"package"`
	Url     string `yaml:"url" json:"url"`
	Branch  string `yaml:"branch" json:"branch"`
	Tag     string `yaml:"tag" json:"tag"`
	Commit  Commit `yaml:"commit" json:"commit"`
}

type Commit struct {
	ID        string `yaml:"id" json:"id"`
	Date      string `yaml:"date" json:"date"`
	Abbr      string `yaml:"abbr" json:"abbr"`
	Sub       string `yaml:"sub" json:"sub"`
	Committer struct {
		Name  string `yaml:"name" json:"name"`
		Email string `yaml:"email" json:"email"`
	} `yaml:"committer" json:"committer"`
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
