package main

import (
	"strings"
	"unicode/utf8"

	"github.com/bidhan-a/indigo/internal/processor"
	"github.com/bidhan-a/indigo/internal/sink"
	"github.com/bidhan-a/indigo/internal/source"
)

func countChars(word interface{}) interface{} {
	w := word.(string)
	w = strings.TrimSuffix(w, "\n")
	return utf8.RuneCountInString(w)
}

func main() {
	source := source.NewFileSource("input.txt")
	processor := processor.NewMap(countChars, 1)
	sink := sink.NewStdSink(false)

	source.AddProcessor(processor).AddSink(sink)
}
