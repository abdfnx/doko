package stream

import (
	"os"
	"io"
	"errors"
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
