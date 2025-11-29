//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Mode -output=mode_string.go

package model

type Mode int32

const (
	ModeIntelligentInferenceScheduling Mode = iota
	ModePrefillDecodeDisaggregation
)
