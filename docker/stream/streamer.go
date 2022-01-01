package stream

import (
	"os"
	"io"
	"log"
	"time"
	"sync"
	"errors"
	"context"

	"github.com/abdfnx/doko/shared"
	logger "github.com/abdfnx/doko/log"

	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/api/types"
)

var (
	ErrEmptyExecID   = errors.New("emtpy exec id")
	ErrTtySizeIsZero = errors.New("tty size is 0")
)

type Streamer struct {
	In    *In
	Out   *Out
	Err   io.Writer
	isTty bool
}

func New() *Streamer {
	return &Streamer{
		In:  NewIn(os.Stdin),
		Out: NewOut(os.Stdout),
		Err: os.Stderr,
	}
}

func (s *Streamer) SetRawTerminal() (func(), error) {
	if err := s.In.SetRawTerminal(); err != nil {
		return nil, err
	}

	var once sync.Once
	restore := func() {
		once.Do(func() {
			if err := s.In.RestoreTerminal(); err != nil {
				logger.Logger.Errorf("failed to restore terminal: %s\n", err)
			}
		})
	}

	return restore, nil
}

func (s *Streamer) resizeTTY(ctx context.Context, resize ResizeContainer, id string) error {
	h, w := s.Out.GetTtySize()
	if h == 0 && w == 0 {
		return ErrTtySizeIsZero
	}

	options := types.ResizeOptions{
		Height: h,
		Width:  w,
	}

	return resize(ctx, id, options)
}

func (s *Streamer) initTTYSize(ctx context.Context, resize ResizeContainer, id string) {
	if err := s.resizeTty(ctx, resize, id); err != nil {
		go func() {
			shared.Logger.Errorf("failed to resize tty: (%s)\n", err)

			for retry := 0; retry < 5; retry++ {
				time.Sleep(10 * time.Millisecond)

				if err = s.resizeTty(ctx, resize, id); err == nil {
					break
				}
			}

			if err != nil {
				log.Println("failed to resize tty, using default size")
			}
		}()
	}
}

func (s *Streamer) streamIn(restore func(), resp types.HijackedResponse) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		_, err := io.Copy(resp.Conn, s.In)

		restore()

		if _, ok := err.(term.EscapeError); ok {
			return
		}

		if err != nil {
			shared.Logger.Errorf("in stream error: %s", err)
			return
		}

		if err := resp.CloseWrite(); err != nil {
			shared.Logger.Errorf("close response error: %s", err)
		}
	}()

	return done
}
