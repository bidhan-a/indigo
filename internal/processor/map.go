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
	m := &Map{
		f:           f,
		in:          make(chan interface{}),
		out:         make(chan interface{}),
		parallelism: parallelism,
	}
	go m.setup()
	return m
}

// InputChan - input channel for Map processor.
func (m *Map) InputChan() chan<- interface{} {
	return m.in
}

// OutputChan - output channel for Map processor.
func (m *Map) OutputChan() <-chan interface{} {
	return m.out
}

// AddProcessor - add a Processor to Map processor.
func (m *Map) AddProcessor(p stream.Processor) stream.Processor {
	go m.connect(p)
	return p
}

// AddSink - add a Sink to Map processor.
func (m *Map) AddSink(s stream.Sink) {
	m.connect(s)
}

func (m *Map) connect(nextProcessorInput stream.Input) {
	for d := range m.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (m *Map) setup() {
	var wg sync.WaitGroup

	for i := 0; i < int(m.parallelism); i++ {
		wg.Add(1)
		go func() {
			for d := range m.in {
				res := m.f(d)
				m.out <- res
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(m.out)
	}()

}
