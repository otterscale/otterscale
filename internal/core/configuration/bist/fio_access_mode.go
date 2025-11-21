package bist

//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=FIOAccessMode -output=fio_access_mode_string.go -linecomment=true

type FIOAccessMode int32

const (
	FIOAccessModeRead          FIOAccessMode = iota // read
	FIOAccessModeWrite                              // write
	FIOAccessModeTrim                               // trim
	FIOAccessModeReadWrite                          // readwrite
	FIOAccessModeTrimWrite                          // trimwrite
	FIOAccessModeRandRead                           // randread
	FIOAccessModeRandWrite                          // randwrite
	FIOAccessModeRandTrim                           // randtrim
	FIOAccessModeRandReadWrite                      // randrw
	FIOAccessModeRandTrimWrite                      // randtrimwrite
)
