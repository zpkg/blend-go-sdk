package async

// RunToError runs all of the functions in separate goroutines and returns the error of the first one to exit
func RunToError(fns ...func() error) error {
	panicChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	for _, fn := range fns {
		go func(fn func() error) {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			errChan <- fn()
		}(fn)
	}

	select {
	case p := <-panicChan:
		panic(p)
	case err := <-errChan:
		return err
	}
}
