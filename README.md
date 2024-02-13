# Go VError

*Inspired by [https://github.com/TritonDataCenter/node-verror](https://github.com/TritonDataCenter/node-verror)*

This module provides a way to wrap errors with extra information as they are returned up the call stack. The `VError` struct implements the `error` interface and supports unwrapping as introduced in go1.13 -- [https://go.dev/blog/go1.13-errors](https://go.dev/blog/go1.13-errors).

## Installation

To use `go-verror`, either install directly with:

```sh
go get -u github.com/jacobtie/go-verror
```

or import in your source code with:

```go
import "github.com/jacobtie/go-verror/verror"
```

and then run:

```sh
go mod tidy
```

## Usage

There are two concepts with `VError` errors. The first is creating an error to optionally wrap an additional error and provide additional information pertaining to the error context. The second is extracting the wrapped information to action off of and/or to log.

### Creating a VError

This module introduces three constructors for creating a VError.

> [!TIP]
> All VError constructors support formatting directives and additional parameters

#### verror.New(msg string)

Creates a new `VError` struct which acts as a drop-in replacement for a normal Go error.

```go
if !isValid(foobar) {
    return verror.New("foobar is invalid")
}
```

As with all `VError` constructors, you may include formatting directives with parameters as well.

```go
if !isValid(foobar) {
    return verror.New("foobar %s is invalid", foobar.Name)
}
```

#### verror.NewWithCause

Creates a new `VError` struct which wraps an error with an additional error string. If wrapping another `VError`, the resulting error will retain the context of the wrapped `VError`.

```go
if err := validate(foobar); err != nil {
    return verror.NewWithCause(err, "foobar is invalid")
}
```

Say that calling `err.Error()` on the error returned by the `validate` function above returns the string `"invalid property baz"`. The resulting `VError` returned by this code snippet would have the error string `"foobar is invalid: invalid property baz"`.

#### verror.NewWithOpts

Creates a new `VError` struct which can optionally wrap an error with an additional error string and optionally include additional information pertaining to the error context.

```go
if err := validate(foobar); err != nil {
    return verror.NewWithOpts(&verror.Options{
        Cause: err,
        Info: map[string]any{
            "foo": "bar",
            "baz": 1,
        },
    }, "foobar %s is invalid", foobar.Name)
}
```

Similar to calling `verror.NewWithCause`, calling `verror.NewWithOpts` will wrap the error message if you provide a `Cause` option. If you do not provide a `Cause` option, then the resulting `VError` will only contain the passed in error string.

Passing in `Info` is where the real power of `VError` comes in and what differentiates `VError` from the built-in error wrapping that comes with `fmt.Errorf` and its `%w` directive. When a `VError` wraps another `VError`, the `Info` is combined. This allows you to build up `Info` as the error is returned up the callstack until you are ready to handle it or simply log it.

### Accessing VError Info

To access the info off a `VError`, call the `verror.Info` function.

```go
info, ok := verror.Info(err)
```

This function will return the shallow copied `info` from the `VError` as a `map[string]any` and an `ok` Boolean. If `ok` is `true`, then `info` will be a non-nil map of wrapped error information. If `ok` is `false`, then the error passed in was `nil` or was not a `VError`, in which case, `info` will be a `nil` map.

A recommended pattern is to continue to return errors wrapped in `VError` structs up the call stack, adding information as you go, until the error reaches error handling code/middleware. At this point, you can use `verror.Info` to extract the error data and, if there is info, to log the information for richer error logs.
