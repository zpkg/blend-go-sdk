package exception

import (
	"encoding/json"
	"fmt"
	"io"
)

var (
	_ error          = (*Ex)(nil)
	_ fmt.Formatter  = (*Ex)(nil)
	_ json.Marshaler = (*Ex)(nil)
)

// New returns a new exception with a call stack.
func New(class interface{}, options ...Option) *Ex {
	return NewWithStackDepth(class, defaultNewStartDepth, options...)
}

// NewWithStackDepth creates a new exception with a given start point of the stack.
func NewWithStackDepth(class interface{}, startDepth int, options ...Option) *Ex {
	if class == nil {
		return nil
	}

	if typed, isTyped := class.(*Ex); isTyped {
		return typed
	} else if err, ok := class.(error); ok {
		return &Ex{
			Class: err,
			Stack: callers(startDepth),
		}
	} else if str, ok := class.(string); ok {
		return &Ex{
			Class: Class(str),
			Stack: callers(startDepth),
		}
	}
	return &Ex{
		Class: Class(fmt.Sprint(class)),
		Stack: callers(startDepth),
	}
}

// Ex is an error with a stack trace.
// It also can have an optional cause, it implements `Exception`
type Ex struct {
	// Class disambiguates between errors, it can be used to identify the type of the error.
	Class error
	// Message adds further detail to the error, and shouldn't be used for disambiguation.
	Message string
	// Inner holds the original error in cases where we're wrapping an error with a stack trace.
	Inner error
	// Stack is the call stack frames used to create the stack output.
	Stack StackTrace
}

// Format allows for conditional expansion in printf statements
// based on the token and flags used.
// 	%+v : class + message + stack
// 	%v, %c : class
// 	%m : message
// 	%t : stack
func (e *Ex) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if e.class != nil && len(e.class.Error()) > 0 {
				fmt.Fprintf(s, "%s", e.class)
				if len(e.message) > 0 {
					fmt.Fprintf(s, "\nmessage: %s", e.message)
				}
			} else if len(e.message) > 0 {
				io.WriteString(s, e.message)
			}
			e.stack.Format(s, verb)
		} else if s.Flag('-') {
			e.stack.Format(s, verb)
		} else {
			io.WriteString(s, e.class.Error())
			if len(e.message) > 0 {
				fmt.Fprintf(s, "\nmessage: %s", e.message)
			}
		}
		if e.inner != nil {
			if typed, ok := e.inner.(fmt.Formatter); ok {
				fmt.Fprint(s, "\ninner: ")
				typed.Format(s, verb)
			} else {
				fmt.Fprintf(s, "\ninner: %v", e.inner)
			}
		}
		return
	case 'c':
		io.WriteString(s, e.class.Error())
	case 'i':
		if e.inner != nil {
			io.WriteString(s, e.inner.Error())
		}
	case 'm':
		io.WriteString(s, e.message)
	case 'q':
		fmt.Fprintf(s, "%q", e.message)
	}
}

// Error implements the `error` interface.
// It returns the exception class, without any of the other supporting context like the stack trace.
// To fetch the stack trace, use .String().
func (e *Ex) Error() string {
	return e.Class.Error()
}

// Decompose breaks the exception down to be marshalled into an intermediate format.
func (e *Ex) Decompose() map[string]interface{} {
	values := map[string]interface{}{}
	values["Class"] = e.Class.Error()
	values["Message"] = e.Message
	if e.stack != nil {
		values["Stack"] = e.Stack.Strings()
	}
	if e.inner != nil {
		if typed, isTyped := e.Inner.(*Ex); isTyped {
			values["Inner"] = typed.Decompose()
		} else {
			values["Inner"] = e.Inner.Error()
		}
	}
	return values
}

// MarshalJSON is a custom json marshaler.
func (e *Ex) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Decompose())
}
