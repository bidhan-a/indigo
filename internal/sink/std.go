package sink

import (
	"fmt"
	"io"
	"os"
)

// StdSink - sink for printing to stdout or stderr.
type StdSink struct {
	writer io.Writer
	in     chan interface{}
}

// NewStdSink - create and return a new StdSink.
func NewStdSink(stderr bool) *StdSink {
	s := &StdSink{
		in: make(chan interface{}),
	}
	if stderr {
		s.writer = os.Stderr
	} else {
		s.writer = os.Stdout
	}
	go s.setup()
	return s
}

// InputChan - input channel for StdSink.
func (s *StdSink) InputChan() chan<- interface{} {
	return s.in
}

func (s *StdSink) setup() {
	for d := range s.in {
		fmt.Fprintln(s.writer, d)
	}
}
