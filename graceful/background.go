/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

import (
	"context"
	"os/signal"
)

// Background yields a context that will signal `<-ctx.Done()` when
// a signal is sent to the process (as specified in `DefaultShutdownSignals`).
func Background() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := Notify(DefaultShutdownSignals...)
	go func() {
		<-shutdown
		cancel()
		signal.Stop(shutdown) // unhook the process signal redirects, the next ^c will crash the process etc.
	}()
	return ctx
}
