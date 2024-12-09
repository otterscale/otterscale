package connector

type Kind int

const (
	KindUnspecified Kind = iota
	KindSource
	KindDestination
)
