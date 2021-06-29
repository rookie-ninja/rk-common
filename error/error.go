package rkerror

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"net/http"
)

type ErrorResp struct {
	Err *Error `json:"error" yaml:"error"`
}

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

type Option func(*ErrorResp)

func FromError(err error) *ErrorResp {
	return &ErrorResp{
		Err: &Error{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: make([]interface{}, 0),
			Message: err.Error(),
		},
	}
}

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

func WithHttpCode(code int) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = code
		resp.Err.Status = http.StatusText(code)
	}
}

func WithGrpcCode(code codes.Code) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = int(code)
		resp.Err.Status = code.String()
	}
}

func WithCodeAndStatus(code int, status string) Option {
	return func(resp *ErrorResp) {
		resp.Err.Code = code
		resp.Err.Status = status
	}
}

func WithMessage(message string) Option {
	return func(resp *ErrorResp) {
		resp.Err.Message = message
	}
}

type Error struct {
	Code    int           `json:"code" yaml:"code"`
	Status  string        `json:"status" yaml:"status"`
	Message string        `json:"message" yaml:"message"`
	Details []interface{} `json:"details" yaml:"details"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%s] %s", err.Status, err.Message)
}
