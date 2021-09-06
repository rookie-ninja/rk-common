// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package rkerror defines RK style API errors.
package rkerror

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"net/http"
)

// ErrorResp is standard rk style error
type ErrorResp struct {
	Err *Error `json:"error" yaml:"error"` // Err is RK style error type
}

// New a error response with options
func New(opts ...Option) *ErrorResp {
	resp := &ErrorResp{
		Err: &Error{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: make([]interface{}, 0),
		},
	}

	for i := range opts {
		opts[i](resp)
	}
	return resp
}

// FromError converts error to ErrorResp
func FromError(err error) *ErrorResp {
	if err == nil {
		err = errors.New("unknown error")
	}

	return &ErrorResp{
		Err: &Error{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: make([]interface{}, 0),
			Message: err.Error(),
		},
	}
}

// Option is ErrorResp option
type Option func(*ErrorResp)

// WithDetails provides any type of error details into error response
func WithDetails(details ...interface{}) Option {
	return func(resp *ErrorResp) {
		for i := range details {
			detail := details[i]

			switch v := detail.(type) {
			case *gin.Error:
				resp.Err.Details = append(resp.Err.Details, v.JSON())
			case *Error:
				resp.Err.Details = append(resp.Err.Details, v.Details...)
			case error:
				resp.Err.Details = append(resp.Err.Details, v.Error())
			default:
				resp.Err.Details = append(resp.Err.Details, v)
			}
		}
	}
}

// WithHttpCode provides response code
func WithHttpCode(code int) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = code
		resp.Err.Status = http.StatusText(code)
	}
}

// WithHttpCode provides grpc response code
func WithGrpcCode(code codes.Code) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = int(code)
		resp.Err.Status = code.String()
	}
}

// WithCodeAndStatus provides http response code and status
func WithCodeAndStatus(code int, status string) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = code
		resp.Err.Status = status
	}
}

// WithMessage provides messages along with response
func WithMessage(message string) Option {
	return func(resp *ErrorResp) {
		resp.Err.Message = message
	}
}

// Error defines standard error types of rk style
type Error struct {
	Code    int           `json:"code" yaml:"code"`       // Code represent codes in response
	Status  string        `json:"status" yaml:"status"`   // Status represent string value of code
	Message string        `json:"message" yaml:"message"` // Message represent detail message
	Details []interface{} `json:"details" yaml:"details"` // Details is a list of details in any types in string
}

// Error returns string of error
func (err *Error) Error() string {
	return fmt.Sprintf("[%s] %s", err.Status, err.Message)
}
