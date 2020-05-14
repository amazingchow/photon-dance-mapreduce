package main

/*
	a word-count application "plugin" for MapReduce.

	go build -buildmode=plugin word_count.go
*/

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/amazingchow/mapreduce/backend/storage"
)

func Map(filename string, contents string) []storage.KeyValue { // nolint
	ff := func(r rune) bool { return !unicode.IsLetter(r) }
	words := strings.FieldsFunc(contents, ff)

	kvs := []storage.KeyValue{}
	for _, w := range words {
		kv := storage.KeyValue{Key: w, Value: "1"}
		kvs = append(kvs, kv)
	}

	return kvs
}

func Reduce(key string, values []string) string { // nolint
	return strconv.Itoa(len(values))
}
