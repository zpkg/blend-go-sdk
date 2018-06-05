exception
=========

This is a simple library for wrapping an `error` with a stack trace.

##Key Concepts

An exception is an error with additional context; message and most importantly, the stack trace at creation.

Concepts:
- `Class`: A grouping error; you should be able to test if exceptions are similar by comparing the exception classes. When an exception is created from an error, the class is set to the original error.
- `Message`: Additional context that is variable, that would otherwise break equatibility with exceptions. You should put extra descriptive information in the message.
- `Inner`: A causing exception or error; if you have to chan multiple errors together as a larger grouped exception chain, use `WithInner(...)`.
- `StackTrace`: A stack of function pointer / frames giving important context to where an exception was created.

##Sample Output

If we run `ex.Error()` on an Exception we will get a more detailed output than a normal `errorString`

```text
Exception: this is a sample error
       At: foo_controller.go:20 testExceptions()
           http.go:198 func1()
           http.go:213 func1()
           http.go:117 func1()
           router.go:299 ServeHTTP()
           server.go:1862 ServeHTTP()
           server.go:1361 serve()
           asm_amd64.s:1696 goexit()
```

##Usage

If we want to create a new exception we can use `New`

```go
	return exception.New("this is a test exception")
```

`New` will create a stack trace at the given line. It ignores stack frames within the `exception` package itself. 

There is also a convenience method `Newf` that will mimic `Sprintf` like behavior.

```go
	return exception.Newf("zone exception: %s", "my zone")
```

Important usage note; to make exceptions more usable you should where possible keep the `Class` of the exception consistent so that you can compare it later.

If you'd like to add variable context to an exception, you can use `WithMessagef(...)`:



If we want to wrap an existing golang `error` all we have to do is call `Wrap`

```go
	file, fileErr := os.ReadFile("my_file.txt")
	if fileErr != nil {
		return exception.Wrap(fileErr)
	}
```

A couple properties of wrap:
* It will return nil if the input error is nil.
* It will not modify an error that is actually an exception, it will simply return it untouched.
* It will create a stack trace for the error if it is not nil, and assign the exception class from the existing error.