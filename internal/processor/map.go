package processor

import (
	"sync"

	"github.com/bidhan-a/indigo/internal/stream"
)

// MapFunc is the 1:1 mapping function.
type MapFunc func(interface{}) interface{}

// Map processor.
type Map struct {
	f           MapFunc
	in          chan interface{}
	out         chan interface{}
	parallelism uint
}

var _ stream.Processor = (*Map)(nil)

// NewMap - initialize and return a new Map processor.
func NewMap(f MapFunc, parallelism uint) *Map {
	mp := &Map{
		f:           f,
		in:          make(chan interface{}),
		out:         make(chan interface{}),
		parallelism: parallelism,
	}
	go mp.setup()
	return mp
}

// InputChan - input channel for Map processor.
func (mp *Map) InputChan() chan<- interface{} {
	return mp.in
}

// OutputChan - output channel for Map processor.
func (mp *Map) OutputChan() <-chan interface{} {
	return mp.out
}

// AddProcessor - add a Processor to Map processor.
func (mp *Map) AddProcessor(p stream.Processor) stream.Processor {
	go mp.connect(p)
	return p
}

// AddSink - add a Sink to Map processor.
func (mp *Map) AddSink(s stream.Sink) {
	mp.connect(s)
}

func (mp *Map) connect(nextProcessorInput stream.Input) {
	for d := range mp.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (mp *Map) setup() {
	var wg sync.WaitGroup

	for i := 0; i < int(mp.parallelism); i++ {
		wg.Add(1)
		go func() {
			for d := range mp.in {
				res := mp.f(d)
				mp.out <- res
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(mp.out)
	}()

}
