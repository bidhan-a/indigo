package main

import (
	"strings"
	"unicode/utf8"

	"github.com/bidhan-a/indigo/internal/processor"
	"github.com/bidhan-a/indigo/internal/sink"
	"github.com/bidhan-a/indigo/internal/source"
)

func min5Chars(word interface{}) bool {
	w := word.(string)
	w = strings.TrimSuffix(w, "\n")
	count := utf8.RuneCountInString(w)
	return count >= 5
}

func toUpper(word interface{}) interface{} {
	w := strings.TrimSuffix(word.(string), "\n")
	return strings.ToUpper(w)
}

func main() {
	source := source.NewFileSource("input.txt")
	mapProcessor := processor.NewMap(toUpper, 1)
	filterProcessor := processor.NewFilter(min5Chars, 1)
	sink := sink.NewStdSink(false)

	source.AddProcessor(filterProcessor).AddProcessor(mapProcessor).AddSink(sink)
}
