package parser

import "go-ip-proxy/result"

type Parser interface {
	Next() bool
	Name() string
	Collect(chan<- *result.Result) []error
}

type Type uint8

const (
	COLLECTBYSELECTOR Type = iota
	COLLECTBYREGEX
)