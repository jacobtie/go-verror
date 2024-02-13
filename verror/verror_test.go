package verror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	err := New("error message")

	assert.Error(err)
	assert.Equal("error message", err.msg)
	assert.Equal("error message", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Nil(err.cause)
	assert.Nil(err.Unwrap())
}

func TestNewFormat(t *testing.T) {
	assert := assert.New(t)

	err := New("error %d message %s", 17, "msg")

	assert.Error(err)
	assert.Equal("error 17 message msg", err.msg)
	assert.Equal("error 17 message msg", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Nil(err.cause)
	assert.Nil(err.Unwrap())
}

func TestNewWithCause(t *testing.T) {
	assert := assert.New(t)

	originalErr := fmt.Errorf("original error")

	err := NewWithCause(originalErr, "new %d error", 17)
	assert.Error(err)
	assert.Equal("new 17 error: original error", err.msg)
	assert.Equal("new 17 error: original error", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Equal(originalErr, err.cause)
	assert.Equal(originalErr, err.Unwrap())
}

func TestNewWithCauseNested(t *testing.T) {
	assert := assert.New(t)

	originalErr1 := fmt.Errorf("original error 1")
	originalErr2 := fmt.Errorf("original error 2: %w", originalErr1)

	err := NewWithCause(originalErr2, "new error")

	assert.Error(err)
	assert.Equal("new error: original error 2: original error 1", err.msg)
	assert.Equal("new error: original error 2: original error 1", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Equal(originalErr2, err.cause)
	assert.Equal(originalErr2, err.Unwrap())
}

func TestNewWithCauseWithInfo(t *testing.T) {
	assert := assert.New(t)

	originalErr := NewWithOpts(&Options{
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "original error")

	err := NewWithCause(originalErr, "new error")

	assert.Error(err)
	assert.Equal("new error: original error", err.msg)
	assert.Equal("new error: original error", err.Error())
	assert.Equal(map[string]interface{}{
		"foo": "bar",
		"baz": 1,
	}, err.info)
	assert.Equal(originalErr, err.cause)
	assert.Equal(originalErr, err.Unwrap())
}

func TestNewWithNestedCauseWithInfo(t *testing.T) {
	assert := assert.New(t)

	originalErr1 := NewWithOpts(&Options{
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "original error 1")
	originalErr2 := NewWithOpts(&Options{
		Cause: originalErr1,
		Info: map[string]interface{}{
			"que": false,
			"baz": 2,
		},
	}, "original error 2")

	err := NewWithCause(originalErr2, "new %d error %s", 17, "msg")

	assert.Error(err)
	assert.Equal("new 17 error msg: original error 2: original error 1", err.msg)
	assert.Equal("new 17 error msg: original error 2: original error 1", err.Error())
	assert.Equal(map[string]interface{}{
		"foo": "bar",
		"baz": 2,
		"que": false,
	}, err.info)
	assert.Equal(originalErr2, err.cause)
	assert.Equal(originalErr2, err.Unwrap())
}

func TestNewWithOptions(t *testing.T) {
	assert := assert.New(t)

	err := NewWithOpts(&Options{
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "error message")

	assert.Error(err)
	assert.Equal("error message", err.msg)
	assert.Equal("error message", err.Error())
	assert.Equal(map[string]interface{}{
		"foo": "bar",
		"baz": 1,
	}, err.info)
	assert.Nil(err.cause)
	assert.Nil(err.Unwrap())
}

func TestNewOptionsWithCause(t *testing.T) {
	assert := assert.New(t)

	originalError := fmt.Errorf("original error")

	err := NewWithOpts(&Options{
		Cause: originalError,
	}, "new error")

	assert.Error(err)
	assert.Equal("new error: original error", err.msg)
	assert.Equal("new error: original error", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Equal(originalError, err.cause)
	assert.Equal(originalError, err.Unwrap())
}

func TestNewOptionsWithNestedCause(t *testing.T) {
	assert := assert.New(t)

	originalErr1 := fmt.Errorf("original error 1")
	originalErr2 := fmt.Errorf("original error 2: %w", originalErr1)

	err := NewWithOpts(&Options{
		Cause: originalErr2,
	}, "new %d error %s", 17, "msg")

	assert.Error(err)
	assert.Equal("new 17 error msg: original error 2: original error 1", err.msg)
	assert.Equal("new 17 error msg: original error 2: original error 1", err.Error())
	assert.Equal(make(map[string]interface{}), err.info)
	assert.Equal(originalErr2, err.cause)
	assert.Equal(originalErr2, err.Unwrap())
}

func TestNewOptionsWithCauseAndInfo(t *testing.T) {
	assert := assert.New(t)

	originalError := fmt.Errorf("original error")

	err := NewWithOpts(&Options{
		Cause: originalError,
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "new error")

	assert.Error(err)
	assert.Equal("new error: original error", err.msg)
	assert.Equal("new error: original error", err.Error())
	assert.Equal(map[string]interface{}{
		"foo": "bar",
		"baz": 1,
	}, err.info)
	assert.Equal(originalError, err.cause)
	assert.Equal(originalError, err.Unwrap())
}

func TestNewOptionsWithCauseAndInfoNested(t *testing.T) {
	assert := assert.New(t)

	originalErr1 := fmt.Errorf("original error 1")
	originalErr2 := NewWithOpts(&Options{
		Cause: originalErr1,
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "original error 2")
	originalErr3 := NewWithOpts(&Options{
		Cause: originalErr2,
		Info: map[string]interface{}{
			"que": false,
		},
	}, "original error 3")

	err := NewWithOpts(&Options{
		Cause: originalErr3,
		Info: map[string]interface{}{
			"baz":    2,
			"foobar": map[string]int{"a": 1, "b": 2},
		},
	}, "new error")

	assert.Error(err)
	assert.Equal("new error: original error 3: original error 2: original error 1", err.msg)
	assert.Equal("new error: original error 3: original error 2: original error 1", err.Error())
	assert.Equal(map[string]interface{}{
		"foo":    "bar",
		"baz":    2,
		"que":    false,
		"foobar": map[string]int{"a": 1, "b": 2},
	}, err.info)
	assert.Equal(originalErr3, err.cause)
	assert.Equal(originalErr3, err.Unwrap())
}

func TestInfo(t *testing.T) {
	assert := assert.New(t)

	originalErr1 := fmt.Errorf("original error 1")
	originalErr2 := NewWithOpts(&Options{
		Cause: originalErr1,
		Info: map[string]interface{}{
			"foo": "bar",
			"baz": 1,
		},
	}, "original error 2")
	originalErr3 := NewWithOpts(&Options{
		Cause: originalErr2,
		Info: map[string]interface{}{
			"que": false,
		},
	}, "original error 3")
	err := NewWithOpts(&Options{
		Cause: originalErr3,
		Info: map[string]interface{}{
			"baz":    2,
			"foobar": map[string]int{"a": 1, "b": 2},
		},
	}, "new error")

	info, ok := Info(err)

	assert.True(ok)
	assert.Equal(map[string]interface{}{
		"foo":    "bar",
		"baz":    2,
		"que":    false,
		"foobar": map[string]int{"a": 1, "b": 2},
	}, info)
}

func TestInfoNotVError(t *testing.T) {
	assert := assert.New(t)

	err := fmt.Errorf("original error 1")
	info, ok := Info(err)

	assert.False(ok)
	assert.Nil(info)
}

func TestUnwrap(t *testing.T) {
	assert := assert.New(t)

	originalErr := fmt.Errorf("original error")
	err := NewWithCause(originalErr, "new error")

	resultErr := Unwrap(err)

	assert.Equal(originalErr, resultErr)
}

func TestUnwrapNil(t *testing.T) {
	assert := assert.New(t)

	err := New("error message")

	resultErr := Unwrap(err)

	assert.Nil(resultErr)
}

func TestUnwrapNotVError(t *testing.T) {
	assert := assert.New(t)

	err := fmt.Errorf("error message")

	resultErr := Unwrap(err)

	assert.Nil(resultErr)
}
