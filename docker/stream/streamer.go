package stream

import (
	"os"
	"io"
	"sync"
	"errors"

	logger "github.com/abdfnx/doko/log"
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
