package stream

// Input ...
type Input interface {
	InputChan() chan<- interface{}
}

// Output ...
type Output interface {
	OutputChan() <-chan interface{}
}

// Source ...
type Source interface {
	Output
	AddProcessor(Processor) Processor
}

// Processor ...
type Processor interface {
	Input
	Output
	AddProcessor(Processor) Processor
	AddSink(Sink)
}

// Sink ...
type Sink interface {
	Input
}
