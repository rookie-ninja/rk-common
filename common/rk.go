// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

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
