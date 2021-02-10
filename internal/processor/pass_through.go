package processor

import (
	"github.com/bidhan-a/indigo/internal/stream"
)

// PassThrough processor.
type PassThrough struct {
	in  chan interface{}
	out chan interface{}
}

var _ stream.Processor = (*PassThrough)(nil)

// NewPassThrough - initialize and return a new PassThrough processor.
func NewPassThrough() *PassThrough {
	ptp := &PassThrough{
		in:  make(chan interface{}),
		out: make(chan interface{}),
	}
	go ptp.setup()
	return ptp
}

// InputChan - input channel for PassThrough processor.
func (ptp *PassThrough) InputChan() chan<- interface{} {
	return ptp.in
}

// OutputChan - output channel for PassThrough processor.
func (ptp *PassThrough) OutputChan() <-chan interface{} {
	return ptp.out
}

// AddProcessor - add a Processor to PassThrough processor.
func (ptp *PassThrough) AddProcessor(p stream.Processor) stream.Processor {
	go ptp.connect(p)
	return p
}

// AddSink - add a Sink to PassThrough processor.
func (ptp *PassThrough) AddSink(s stream.Sink) {
	ptp.connect(s)
}

func (ptp *PassThrough) connect(nextProcessorInput stream.Input) {
	for d := range ptp.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (ptp *PassThrough) setup() {
	for d := range ptp.in {
		ptp.out <- d
	}
	close(ptp.out)

}
