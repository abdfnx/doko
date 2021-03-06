//go:build !windows
// +build !windows

package stream

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func (s *Streamer) monitorTtySize(ctx context.Context, resize ResizeContainer, id string) {
	s.initTTYSize(ctx, resize, id)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGWINCH)
	go func() {
		for range sigchan {
			s.resizeTTY(ctx, resize, id)
		}
	}()
}
