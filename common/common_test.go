// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkcommon

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestGetDefaultIfEmptyString_ExpectDefault(t *testing.T) {
	def := "unit-test-default"
	assert.Equal(t, def, GetDefaultIfEmptyString("", def))
}

func TestGetDefaultIfEmptyString_ExpectOriginal(t *testing.T) {
	def := "unit-test-default"
	origin := "init-test-original"
	assert.Equal(t, origin, GetDefaultIfEmptyString(origin, def))
}

func TestFileExists_ExpectTrue(t *testing.T) {
	filePath := path.Join(t.TempDir(), "ui-TestFileExist-ExpectTrue")
	assert.Nil(t, ioutil.WriteFile(filePath, []byte("unit-test"), 0777))
	assert.True(t, FileExists(filePath))
}

func TestFileExists_ExpectFalse(t *testing.T) {
	filePath := path.Join(t.TempDir(), "ui-TestFileExist-ExpectFalse")
	assert.False(t, FileExists(filePath))
}

func TestFileExists_WithEmptyFilePath(t *testing.T) {
	assert.False(t, FileExists(""))
}

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

func TestGetEnvValueOrDefault_ExpectEnvValue(t *testing.T) {
	assert.Nil(t, os.Setenv("unit-test-key", "unit-test-value"))
	assert.Equal(t, "unit-test-value", GetEnvValueOrDefault("unit-test-key", "unit-test-default"))
	os.Clearenv()
}

func TestGetEnvValueOrDefault_ExpectDefaultValue(t *testing.T) {
	assert.Equal(t, "unit-test-default", GetEnvValueOrDefault("unit-test-key", "unit-test-default"))
}

func TestGetLocalIP_HappyCase(t *testing.T) {
	assert.NotEmpty(t, GetLocalIP())
}

func TestGetLocalHostname_HappyCase(t *testing.T) {
	assert.NotEmpty(t, GetLocalHostname())
}

func TestGenerateRequestId_HappyCase(t *testing.T) {
	assert.NotEmpty(t, GenerateRequestId())
}

func TestGenerateRequestIdWithPrefix_HappyCase(t *testing.T) {
	assert.True(t, strings.HasPrefix(GenerateRequestIdWithPrefix("unit-test"), "unit-test"))
}

func TestOverrideMap_WithNilSource(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	override := make(map[interface{}]interface{})
	OverrideMap(nil, override)
}

func TestOverrideMap_WithNilOverride(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := make(map[interface{}]interface{})
	OverrideMap(src, nil)
}

func TestOverrideMap_WithUnMatchedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := make(map[interface{}]interface{})
	src["ut-src-key"] = "ut-src-value"

	override := make(map[interface{}]interface{})
	override["ut-override-key"] = false

	OverrideMap(src, override)

	// source map should keep the same
	assert.Equal(t, 1, len(src))
	assert.Equal(t, "ut-src-value", src["ut-src-key"])

	// override map should never change
	assert.Equal(t, 1, len(override))
	assert.Equal(t, false, override["ut-override-key"])
}

func TestOverrideMap_WithMixedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := make(map[interface{}]interface{})
	src["ut-src-key"] = "ut-src-value"

	override := make(map[interface{}]interface{})
	override["ut-override-key"] = false
	override["ut-src-key"] = "ut-override-value"

	OverrideMap(src, override)

	// source map should be changed
	assert.Equal(t, 1, len(src))
	assert.Equal(t, "ut-override-value", src["ut-src-key"])

	// override map should never change
	assert.Equal(t, 2, len(override))
	assert.Equal(t, false, override["ut-override-key"])
	assert.Equal(t, "ut-override-value", src["ut-src-key"])
}

func TestOverrideMap_WithHappyCase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	type MyStruct struct {
		Key string
	}

	// no panic expected
	src := make(map[interface{}]interface{})
	src["ut-str"] = ""
	src["ut-list"] = []string{}
	src["ut-map"] = map[string]interface{}{}
	src["ut-struct"] = &MyStruct{}

	override := make(map[interface{}]interface{})
	override["ut-str"] = "ut-str"
	override["ut-list"] = []string{"one", "two"}
	override["ut-map"] = map[string]interface{}{
		"key": "value",
	}
	override["ut-struct"] = &MyStruct{
		Key: "override",
	}

	OverrideMap(src, override)

	// source map should be changed
	// validate string
	assert.Equal(t, "ut-str", src["ut-str"])
	// validate list
	assert.Contains(t, src["ut-list"], "one")
	assert.Contains(t, src["ut-list"], "two")
	// validate map
	innerMap := src["ut-map"]
	assert.Equal(t, "value", innerMap.(map[string]interface{})["key"])
	// validate struct
	innerStruct := src["ut-struct"]
	assert.NotNil(t, "override", innerStruct.(*MyStruct).Key)
}

func TestOverrideSlice_WithNilSource(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	override := make([]interface{}, 0)
	OverrideSlice(nil, override)
}

func TestOverrideSlice_WithNilOverride(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := make([]interface{}, 0)
	OverrideSlice(src, nil)
}

func TestOverrideSlice_WithUnMatchedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := []interface{}{"str"}
	override := []interface{}{false}

	OverrideSlice(src, override)

	assert.Len(t, src, 1)
	assert.Equal(t, "str", src[0])

	assert.Len(t, override, 1)
	assert.Equal(t, false, override[0])
}

func TestOverrideSlice_WithMixedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	// no panic expected
	src := []interface{}{"str", true}
	override := []interface{}{"override-str", false}

	OverrideSlice(src, override)

	assert.Len(t, src, 2)
	assert.Equal(t, "override-str", src[0])
	assert.Equal(t, false, src[1])

	assert.Len(t, override, 2)
	assert.Equal(t, "override-str", src[0])
	assert.Equal(t, false, src[1])
}

func TestOverrideSlice_WithHappyCase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// this should never be called in case of a bug
			assert.True(t, false)
		} else {
			// no panic expected
			assert.True(t, true)
		}
	}()

	type MyStruct struct {
		Key string
	}

	// no panic expected
	src := []interface{}{
		"",
		[]string{},
		map[string]interface{}{},
		&MyStruct{},
	}

	override := []interface{}{
		"str",
		[]string{"one", "two"},
		map[string]interface{}{"key": "value"},
		&MyStruct{Key: "override"},
	}

	OverrideSlice(src, override)

	// source map should be changed
	// validate string
	assert.Equal(t, "str", src[0])
	// validate list
	assert.Contains(t, src[1], "one")
	assert.Contains(t, src[1], "two")
	// validate map
	innerMap := src[2]
	assert.Equal(t, "value", innerMap.(map[string]interface{})["key"])
	// validate struct
	innerStruct := src[3]
	assert.NotNil(t, "override", innerStruct.(*MyStruct).Key)
}

func TestConvertJSONToMap_WithEmptyString(t *testing.T) {
	assert.Empty(t, ConvertJSONToMap(""))
}

func TestConvertJSONToMap_WithInvalidString(t *testing.T) {
	assert.Empty(t, ConvertJSONToMap("{"))
}

func TestConvertJSONToMap_WithInvalidJSON(t *testing.T) {
	assert.Empty(t, ConvertJSONToMap("{]"))
}

func TestConvertJSONToMap_HappyCase(t *testing.T) {
	str := `{"key":"value"}`

	res := ConvertJSONToMap(str)

	assert.Equal(t, "value", res["key"])
}

func TestConvertStructToJSON_WithNilStruct(t *testing.T) {
	assert.Empty(t, ConvertStructToJSON(nil))
}

func TestConvertStructToJSON_WithNonJSONStruct(t *testing.T) {
	type MyStruct struct {
		KEY   string
		MAP   map[string]interface{}
		SLICE []interface{}
	}

	res := ConvertStructToJSON(&MyStruct{
		KEY: "key",
		MAP: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		SLICE: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "KEY")
	assert.Contains(t, res, "MAP")
	assert.Contains(t, res, "SLICE")
}

func TestConvertStructToJSON_HappyCase(t *testing.T) {
	type MyStruct struct {
		Key   string                 `json:"key"`
		Map   map[string]interface{} `json:"map"`
		Slice []interface{}          `json:"slice"`
	}

	res := ConvertStructToJSON(&MyStruct{
		Key: "key",
		Map: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		Slice: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "key")
	assert.Contains(t, res, "map")
	assert.Contains(t, res, "slice")
}

func TestConvertStructToJSONPretty_WithNilStruct(t *testing.T) {
	assert.Empty(t, ConvertStructToJSONPretty(nil))
}

func TestConvertStructToJSONPretty_WithNonJSONStruct(t *testing.T) {
	type MyStruct struct {
		KEY   string
		MAP   map[string]interface{}
		SLICE []interface{}
	}

	res := ConvertStructToJSONPretty(&MyStruct{
		KEY: "key",
		MAP: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		SLICE: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "KEY")
	assert.Contains(t, res, "MAP")
	assert.Contains(t, res, "SLICE")
	assert.Contains(t, res, "  ")
}

func TestConvertStructToJSONPretty_HappyCase(t *testing.T) {
	type MyStruct struct {
		Key   string                 `json:"key"`
		Map   map[string]interface{} `json:"map"`
		Slice []interface{}          `json:"slice"`
	}

	res := ConvertStructToJSONPretty(&MyStruct{
		Key: "key",
		Map: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		Slice: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "key")
	assert.Contains(t, res, "map")
	assert.Contains(t, res, "slice")
	assert.Contains(t, res, "  ")
}

func TestConvertStructToBytes_WithNilStruct(t *testing.T) {
	assert.Empty(t, ConvertStructToBytes(nil))
}

func TestConvertStructToBytes_WithNonJSONStruct(t *testing.T) {
	type MyStruct struct {
		KEY   string
		MAP   map[string]interface{}
		SLICE []interface{}
	}

	res := ConvertStructToBytes(&MyStruct{
		KEY: "key",
		MAP: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		SLICE: []interface{}{
			false,
			"str",
		},
	})

	str := string(res)
	assert.Contains(t, str, "KEY")
	assert.Contains(t, str, "MAP")
	assert.Contains(t, str, "SLICE")
}

func TestConvertStructToBytes_HappyCase(t *testing.T) {
	type MyStruct struct {
		Key   string                 `json:"key"`
		Map   map[string]interface{} `json:"map"`
		Slice []interface{}          `json:"slice"`
	}

	res := ConvertStructToBytes(&MyStruct{
		Key: "key",
		Map: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		Slice: []interface{}{
			false,
			"str",
		},
	})

	str := string(res)
	assert.Contains(t, str, "key")
	assert.Contains(t, str, "map")
	assert.Contains(t, str, "slice")
}

func TestConvertStructToMap_WithNilStruct(t *testing.T) {
	assert.Empty(t, ConvertStructToMap(nil))
}

func TestConvertStructToMap_WithNonJSONStruct(t *testing.T) {
	type MyStruct struct {
		KEY   string
		MAP   map[string]interface{}
		SLICE []interface{}
	}

	res := ConvertStructToMap(&MyStruct{
		KEY: "key",
		MAP: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		SLICE: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "KEY")
	assert.Contains(t, res, "MAP")
	assert.Contains(t, res, "SLICE")
}

func TestConvertStructToMap_HappyCase(t *testing.T) {
	type MyStruct struct {
		Key   string                 `json:"key"`
		Map   map[string]interface{} `json:"map"`
		Slice []interface{}          `json:"slice"`
	}

	res := ConvertStructToMap(&MyStruct{
		Key: "key",
		Map: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		Slice: []interface{}{
			false,
			"str",
		},
	})

	assert.Contains(t, res, "key")
	assert.Contains(t, res, "map")
	assert.Contains(t, res, "slice")
}

func TestConvertStructToZapFields_WithNilStruct(t *testing.T) {
	assert.Empty(t, ConvertStructToZapFields(nil))
}

func TestConvertStructToZapFields_WithNonJSONStruct(t *testing.T) {
	type MyStruct struct {
		KEY   string
		MAP   map[string]interface{}
		SLICE []interface{}
	}

	res := ConvertStructToZapFields(&MyStruct{
		KEY: "key",
		MAP: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		SLICE: []interface{}{
			false,
			"str",
		},
	})

	assert.Len(t, res, 3)
}

func TestConvertStructToZapFields_HappyCase(t *testing.T) {
	type MyStruct struct {
		Key   string                 `json:"key"`
		Map   map[string]interface{} `json:"map"`
		Slice []interface{}          `json:"slice"`
	}

	res := ConvertStructToZapFields(&MyStruct{
		Key: "key",
		Map: map[string]interface{}{
			"bool": false,
			"str":  "str",
		},
		Slice: []interface{}{
			false,
			"str",
		},
	})

	assert.Len(t, res, 3)
}

func TestMatchLocaleWithEnv_WithEmptyLocale(t *testing.T) {
	assert.False(t, MatchLocaleWithEnv(""))
}

func TestMatchLocaleWithEnv_WithInvalidLocale(t *testing.T) {
	assert.False(t, MatchLocaleWithEnv("realm::region::az"))
}

func TestMatchLocaleWithEnv_WithEmptyRealmEnv(t *testing.T) {
	// with realm exist in locale
	assert.False(t, MatchLocaleWithEnv("fake-realm::*::*::*"))

	// with wildcard in realm
	assert.True(t, MatchLocaleWithEnv("*::*::*::*"))
}

func TestMatchLocaleWithEnv_WithRealmEnv(t *testing.T) {
	// set environment variable
	assert.Nil(t, os.Setenv("REALM", "ut"))

	// with realm exist in locale
	assert.True(t, MatchLocaleWithEnv("ut::*::*::*"))

	// with wildcard in realm
	assert.True(t, MatchLocaleWithEnv("*::*::*::*"))

	// with wrong realm
	assert.False(t, MatchLocaleWithEnv("rk::*::*::*"))

	assert.Nil(t, os.Setenv("REALM", ""))
}

func TestMatchLocaleWithEnv_WithRegionEnv(t *testing.T) {
	// set environment variable
	assert.Nil(t, os.Setenv("REGION", "ut"))

	// with region exist in locale
	assert.True(t, MatchLocaleWithEnv("*::ut::*::*"))

	// with wildcard in region
	assert.True(t, MatchLocaleWithEnv("*::*::*::*"))

	// with wrong region
	assert.False(t, MatchLocaleWithEnv("*::rk::*::*"))

	assert.Nil(t, os.Setenv("REGION", ""))
}

func TestMatchLocaleWithEnv_WithAZEnv(t *testing.T) {
	// set environment variable
	assert.Nil(t, os.Setenv("AZ", "ut"))

	// with az exist in locale
	assert.True(t, MatchLocaleWithEnv("*::*::ut::*"))

	// with wildcard in az
	assert.True(t, MatchLocaleWithEnv("*::*::*::*"))

	// with wrong az
	assert.False(t, MatchLocaleWithEnv("*::*::rk::*"))

	assert.Nil(t, os.Setenv("AZ", ""))
}

func TestMatchLocaleWithEnv_WithDomainEnv(t *testing.T) {
	// set environment variable
	assert.Nil(t, os.Setenv("DOMAIN", "ut"))

	// with domain exist in locale
	assert.True(t, MatchLocaleWithEnv("*::*::*::ut"))

	// with wildcard in domain
	assert.True(t, MatchLocaleWithEnv("*::*::*::*"))

	// with wrong domain
	assert.False(t, MatchLocaleWithEnv("*::*::*::rk"))

	assert.Nil(t, os.Setenv("DOMAIN", ""))
}

func TestGetUsernameFromBasicAuthString_WithInvalidBasicAuth(t *testing.T) {
	assert.Empty(t, GetUsernameFromBasicAuthString("invalid-basic-auth"))
}

func TestGetUsernameFromBasicAuthString_HappyCase(t *testing.T) {
	assert.Equal(t, "user", GetUsernameFromBasicAuthString("user:pass"))
}

func TestGetPasswordFromBasicAuthString_WithInvalidBasicAuth(t *testing.T) {
	assert.Empty(t, GetPasswordFromBasicAuthString("invalid-basic-auth"))
}

func TestGetPasswordFromBasicAuthString_HappyCase(t *testing.T) {
	assert.Equal(t, "pass", GetPasswordFromBasicAuthString("user:pass"))
}

func TestExtractSchemeFromURL_InvalidURL(t *testing.T) {
	assert.Empty(t, ExtractSchemeFromURL("ftp://localhost"))
}

func TestExtractSchemeFromURL_WithHTTP(t *testing.T) {
	assert.Equal(t, "http", ExtractSchemeFromURL("http://localhost"))
}

func TestExtractSchemeFromURL_WithHTTPS(t *testing.T) {
	assert.Equal(t, "https", ExtractSchemeFromURL("https://localhost"))
}

func TestGetLocale_WithoutEnvVariables(t *testing.T) {
	assert.Equal(t, "*::*::*::*", GetLocale())
}

func TestGetLocale_WithRealm(t *testing.T) {
	os.Setenv("REALM", "ut-realm")
	defer os.Setenv("REALM", "")
	assert.Equal(t, "ut-realm::*::*::*", GetLocale())
}

func TestGetLocale_WithRegion(t *testing.T) {
	os.Setenv("REGION", "ut-region")
	defer os.Setenv("REGION", "")
	assert.Equal(t, "*::ut-region::*::*", GetLocale())
}

func TestGetLocale_WithAZ(t *testing.T) {
	os.Setenv("AZ", "ut-az")
	defer os.Setenv("AZ", "")
	assert.Equal(t, "*::*::ut-az::*", GetLocale())
}

func TestGetLocale_WithDomain(t *testing.T) {
	os.Setenv("DOMAIN", "ut-domain")
	defer os.Setenv("DOMAIN", "")
	assert.Equal(t, "*::*::*::ut-domain", GetLocale())
}
