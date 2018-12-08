package sh

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// Pipe runs commands in sequence, piping combined output to the standard in of the next command.
func Pipe(commands ...*exec.Cmd) error {
	wg := sync.WaitGroup{}
	wg.Add(len(commands))

	errors := make(chan error, len(commands))
	readers := make([]io.Reader, len(commands))
	writers := make([]io.Writer, len(commands))
	for index := 0; index < len(commands); index++ {
		// set up pipes
		readers[index], writers[index] = io.Pipe()

		// wire up pipes
		if index == 0 { // the first command
			commands[index].Stdin = os.Stdin
			commands[index].Stdout = writers[index]
			commands[index].Stderr = writers[index]
		} else if index == len(commands)-1 { // the last command
			commands[index].Stdin = readers[index-1]
			commands[index].Stdout = os.Stdout
			commands[index].Stderr = os.Stderr
		} else { // intermediate commands
			commands[index].Stdin = readers[index-1]
			commands[index].Stdout = writers[index]
			commands[index].Stderr = writers[index]
		}

		go func(index int, cmd *exec.Cmd) {
			defer wg.Done()
			defer func() {
				if err := readers[index].(*io.PipeReader).Close(); err != nil {
					errors <- err
				}
				writers[index] = nil
			}()
			defer func() {
				if err := writers[index].(*io.PipeWriter).Close(); err != nil {
					errors <- err
				}
				writers[index] = nil
			}()
			if index == 1 {
				defer func() {
					if closer, ok := cmd.Stdin.(io.Closer); ok {
						if err := closer.Close(); err != nil {
							errors <- err
						}
					}
				}()
			}
			if err := cmd.Run(); err != nil {
				if !IsEPIPE(err) {
					errors <- err
				}
			}

		}(index, commands[index])
	}

	wg.Wait()
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}

// IsEPIPE is the epipe erorr.
func IsEPIPE(err error) bool {
	if typed, ok := err.(*exec.ExitError); ok {
		status := typed.Sys().(syscall.WaitStatus)
		if status.Signal() == syscall.SIGPIPE {
			return true
		}
	}
	return false
}
