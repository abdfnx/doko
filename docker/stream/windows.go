//go:build windows
// +build windows

package streamer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func (s *Streamer) monitorTtySize(ctx context.Context, resize ResizeContainer, id string) {
	// TODO: add support for Windows
}
