package stream

import "github.com/docker/docker/pkg/term"

type CommonStream struct {
	State      *term.State
	IsTerminal bool
	Fd         uintptr
}

func (s *CommonStream) RestoreTerminal() {
	if s.State != nil {
		term.RestoreTerminal(s.Fd, s.State)
	}
}
