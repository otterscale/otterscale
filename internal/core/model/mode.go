//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Mode -output=mode_string.go -linecomment=true

package model

type Mode int32

const (
	ModeIntelligentInferenceScheduling Mode = iota // intelligent-inference-scheduling
	ModePrefillDecodeDisaggregation                // prefill-decode-disaggregation
)
