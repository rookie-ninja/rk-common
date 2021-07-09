package rkerror

import (
	"fmt"
	"github.com/rookie-ninja/rk-common/error/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorWrapper func(msg string, errors ...error) *status.Status

var (
	Canceled           = BaseErrorWrapper(codes.Canceled)
	Unknown            = BaseErrorWrapper(codes.Unknown)
	InvalidArgument    = BaseErrorWrapper(codes.InvalidArgument)
	DeadlineExceeded   = BaseErrorWrapper(codes.DeadlineExceeded)
	NotFound           = BaseErrorWrapper(codes.NotFound)
	AlreadyExists      = BaseErrorWrapper(codes.AlreadyExists)
	PermissionDenied   = BaseErrorWrapper(codes.PermissionDenied)
	ResourceExhausted  = BaseErrorWrapper(codes.ResourceExhausted)
	FailedPrecondition = BaseErrorWrapper(codes.FailedPrecondition)
	Aborted            = BaseErrorWrapper(codes.Aborted)
	OutOfRange         = BaseErrorWrapper(codes.OutOfRange)
	Unimplemented      = BaseErrorWrapper(codes.Unimplemented)
	Internal           = BaseErrorWrapper(codes.Internal)
	Unavailable        = BaseErrorWrapper(codes.Unavailable)
	DataLoss           = BaseErrorWrapper(codes.DataLoss)
	Unauthenticated    = BaseErrorWrapper(codes.Unauthenticated)
)

func BaseErrorWrapper(code codes.Code) ErrorWrapper {
	return func(msg string, errors ...error) *status.Status {
		st := status.New(code, msg)

		// Inject grpc error as detail
		st, _ = st.WithDetails(&rk_error.ErrorDetail{
			Code: int32(code),
			Status: code.String(),
			Message: fmt.Sprintf("[from-grpc] %s", msg),
		})

		for i := range errors {
			st1, _ := status.FromError(errors[i])
			detail := &rk_error.ErrorDetail{
				Code:    int32(st1.Code()),
				Status:  st1.Code().String(),
				Message: st1.Message(),
			}
			st, _ = st.WithDetails(detail)
		}

		return st
	}
}
