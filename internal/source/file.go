package source

import (
	"bufio"
	"log"
	"os"

	"github.com/bidhan-a/indigo/internal/stream"
)

// FileSource - file source.
type FileSource struct {
	fileName string
	out      chan interface{}
}

// NewFileSource - create and return a new FileSource.
func NewFileSource(fileName string) *FileSource {
	f := &FileSource{
		fileName: fileName,
		out:      make(chan interface{}),
	}
	go f.setup()
	return f
}

// OutputChan - output channel for FileSource.
func (f *FileSource) OutputChan() <-chan interface{} {
	return f.out
}

// AddProcessor - add a Processor to FileSource.
func (f *FileSource) AddProcessor(p stream.Processor) stream.Processor {
	go f.connect(p)
	return p
}

func (f *FileSource) connect(nextProcessorInput stream.Input) {
	for d := range f.OutputChan() {
		nextProcessorInput.InputChan() <- d
	}
	close(nextProcessorInput.InputChan())
}

func (f *FileSource) setup() {
	file, err := os.Open(f.fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		l, isPrefix, err := reader.ReadLine()
		if err != nil {
			close(f.out)
			break
		}

		var line string
		if isPrefix {
			// incomplete line
			line = string(l)
		} else {
			// complete line
			line = string(l) + "\n"
		}

		f.out <- line
	}
}
