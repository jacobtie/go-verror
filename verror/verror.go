package verror

import "fmt"

type VError struct {
	msg   string
	info  map[string]any
	cause error
}

type Options struct {
	Info  map[string]any
	Cause error
}

func New(msg string, params ...any) *VError {
	return &VError{
		msg:   fmt.Sprintf(msg, params...),
		info:  make(map[string]any),
		cause: nil,
	}
}

func NewWithCause(cause error, msg string, params ...any) *VError {
	errMsg := fmt.Sprintf(msg, params...)
	errInfo := make(map[string]any)
	if cause != nil {
		errMsg = fmt.Sprintf("%s: %s", errMsg, cause.Error())
		if causeInfo, ok := Info(cause); ok {
			errInfo = causeInfo
		}
	}
	return &VError{
		msg:   errMsg,
		info:  errInfo,
		cause: cause,
	}
}

func NewWithOpts(opts *Options, msg string, params ...any) *VError {
	errMsg := fmt.Sprintf(msg, params...)
	errInfo := make(map[string]any)
	var errCause error = nil
	if opts != nil {
		if opts.Cause != nil {
			errMsg = fmt.Sprintf("%s: %s", errMsg, opts.Cause.Error())
			if causeInfo, ok := Info(opts.Cause); ok {
				for k, v := range causeInfo {
					errInfo[k] = v
				}
			}
			errCause = opts.Cause
		}
		for k, v := range opts.Info {
			errInfo[k] = v
		}
	}
	return &VError{
		msg:   errMsg,
		info:  errInfo,
		cause: errCause,
	}
}

func (v *VError) Error() string {
	return v.msg
}

func (v *VError) Unwrap() error {
	return v.cause
}

func Info(err error) (map[string]any, bool) {
	verror, ok := err.(*VError)
	if !ok || verror == nil {
		return nil, false
	}
	errInfo := make(map[string]any)
	for k, v := range verror.info {
		errInfo[k] = v
	}
	return errInfo, true
}

func Unwrap(err error) error {
	errWithUnwrap, ok := err.(interface{ Unwrap() error })
	if !ok || errWithUnwrap == nil {
		return nil
	}
	return errWithUnwrap.Unwrap()
}
