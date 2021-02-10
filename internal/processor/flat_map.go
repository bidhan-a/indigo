package processor

import (
	"sync"

	"github.com/bidhan-a/indigo/internal/stream"
)

// FlatMapFunc is the 1:N mapping function.
type FlatMapFunc func(interface{}) []interface{}

// FlatMap processor.
type FlatMap struct {
	f           FlatMapFunc
	in          chan interface{}
	out         chan interface{}
	parallelism uint
}

var _ stream.Processor = (*FlatMap)(nil)

// NewFlatMap - initialize and return a new Map processor.
func NewFlatMap(f FlatMapFunc, parallelism uint) *FlatMap {
	fmp := &FlatMap{
		f:           f,
		in:          make(chan interface{}),
		out:         make(chan interface{}),
		parallelism: parallelism,
	}
	go fmp.setup()
	return fmp
}

// InputChan - input channel for FlatMap processor.
func (fmp *FlatMap) InputChan() chan<- interface{} {
	return fmp.in
}

// OutputChan - output channel for FlatMap processor.
func (fmp *FlatMap) OutputChan() <-chan interface{} {
	return fmp.out
}

// AddProcessor - add a Processor to FlatMap processor.
func (fmp *FlatMap) AddProcessor(p stream.Processor) stream.Processor {
	go fmp.connect(p)
	return p
}

// AddSink - add a Sink to FlatMap processor.
func (fmp *FlatMap) AddSink(s stream.Sink) {
	fmp.connect(s)
}

func (fmp *FlatMap) connect(nextProcessorInput stream.Input) {
	for d := range fmp.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (fmp *FlatMap) setup() {
	var wg sync.WaitGroup

	for i := 0; i < int(fmp.parallelism); i++ {
		wg.Add(1)
		go func() {
			for d := range fmp.in {
				res := fmp.f(d)
				fmp.out <- res
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(fmp.out)
	}()

}
