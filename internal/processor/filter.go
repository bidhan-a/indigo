package processor

import (
	"sync"

	"github.com/bidhan-a/indigo/internal/stream"
)

// FilterFunc is the filtering function.
type FilterFunc func(interface{}) bool

// Filter processor.
type Filter struct {
	f           FilterFunc
	in          chan interface{}
	out         chan interface{}
	parallelism uint
}

var _ stream.Processor = (*Filter)(nil)

// NewFilter - initialize and return a new Filter processor.
func NewFilter(f FilterFunc, parallelism uint) *Filter {
	fp := &Filter{
		f:           f,
		in:          make(chan interface{}),
		out:         make(chan interface{}),
		parallelism: parallelism,
	}
	go fp.setup()
	return fp
}

// InputChan - input channel for Filter processor.
func (fp *Filter) InputChan() chan<- interface{} {
	return fp.in
}

// OutputChan - output channel for Filter processor.
func (fp *Filter) OutputChan() <-chan interface{} {
	return fp.out
}

// AddProcessor - add a Processor to Filter processor.
func (fp *Filter) AddProcessor(p stream.Processor) stream.Processor {
	go fp.connect(p)
	return p
}

// AddSink - add a Sink to Filter processor.
func (fp *Filter) AddSink(s stream.Sink) {
	fp.connect(s)
}

func (fp *Filter) connect(nextProcessorInput stream.Input) {
	for d := range fp.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (fp *Filter) setup() {
	var wg sync.WaitGroup

	for i := 0; i < int(fp.parallelism); i++ {
		wg.Add(1)
		go func() {
			for d := range fp.in {
				res := fp.f(d)
				if res {
					fp.out <- d
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(fp.out)
	}()

}
