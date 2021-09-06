// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkerror

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"net/http"
	"testing"
)

func TestNew_WithoutOptions(t *testing.T) {
	res := New()

	assert.NotNil(t, res)
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Empty(t, res.Err.Details)
}

func TestNew_WithDetails(t *testing.T) {
	// With gin error
	ginErr := &gin.Error{
		Type: gin.ErrorTypePublic,
		Err:  errors.New("gin error"),
		Meta: "gin meta",
	}

	res := New(WithDetails(ginErr))

	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Equal(t, "gin meta", res.Err.Details[0].(gin.H)["meta"])
	assert.Equal(t, "gin error", res.Err.Details[0].(gin.H)["error"])

	// With rk error type
	res = New(WithDetails("rk error"))
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Equal(t, "rk error", res.Err.Details[0])

	// With go error type
	res = New(WithDetails(errors.New("go error")))
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Equal(t, "go error", res.Err.Details[0])

	// With other type
	res = New(WithDetails("error string"))
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Equal(t, "error string", res.Err.Details[0])
}

func TestNew_WithHttpCode(t *testing.T) {
	res := New(WithHttpCode(http.StatusAlreadyReported))

	assert.Equal(t, http.StatusAlreadyReported, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusAlreadyReported), res.Err.Status)
}

func TestNew_WithGrpcCode(t *testing.T) {
	res := New(WithGrpcCode(codes.Aborted))

	assert.Equal(t, int(codes.Aborted), res.Err.Code)
	assert.Equal(t, codes.Aborted.String(), res.Err.Status)
}

func TestNew_WithCodeAndStatus(t *testing.T) {
	res := New(WithCodeAndStatus(-1, "ut-status"))

	assert.Equal(t, -1, res.Err.Code)
	assert.Equal(t, "ut-status", res.Err.Status)
}

func TestNew_WithMessage(t *testing.T) {
	res := New(WithMessage("ut message"))

	assert.Equal(t, "ut message", res.Err.Message)
}

func TestFromError_WithNilError(t *testing.T) {
	res := FromError(nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Empty(t, res.Err.Details)
	assert.Equal(t, "unknown error", res.Err.Message)
}

func TestFromError_HappyCase(t *testing.T) {
	res := FromError(errors.New("ut error"))
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusInternalServerError, res.Err.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), res.Err.Status)
	assert.Empty(t, res.Err.Details)
	assert.Equal(t, "ut error", res.Err.Message)
}
