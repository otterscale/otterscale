// cmd_test.go
package ceph // same package as the implementation – can access unexported functions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ---------------------------------------------------------- *
 * Helper – a type that JSON cannot marshal (used to provoke an
 * error in the first step of the command helpers).
 * ---------------------------------------------------------- */
type unmarshalable struct {
	Ch chan int // channels are not JSON‑serialisable
}

/* ---------------------------------------------------------- *
 * osdCommand / monCommand – marshal error path
 * ---------------------------------------------------------- */
func TestOSDCommand_MarshalError(t *testing.T) {
	// osdCommand first marshals the command; the channel makes it fail.
	var bad unmarshalable
	// conn is irrelevant because marshal fails before it is used.
	_, err := osdCommand(nil, 0, bad)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

func TestMonCommand_MarshalError(t *testing.T) {
	var bad unmarshalable
	_, err := monCommand(nil, bad)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

/* ---------------------------------------------------------- *
 * osdCommandWithKeyUnmarshal – error propagation when the
 * underlying osdCommand fails (marshal error above)
 * ---------------------------------------------------------- */
func TestOSDCommandWithKeyUnmarshal_PropagatesError(t *testing.T) {
	var bad unmarshalable
	var dst any
	err := osdCommandWithKeyUnmarshal(nil, 0, "anyKey", bad, &dst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

/* ---------------------------------------------------------- *
 * monCommandWithUnmarshal – error propagation when the
 * underlying monCommand fails (marshal error above)
 * ---------------------------------------------------------- */
func TestMonCommandWithUnmarshal_PropagatesError(t *testing.T) {
	var bad unmarshalable
	var dst any
	err := monCommandWithUnmarshal(nil, bad, &dst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

/* ---------------------------------------------------------- *
 * monCommandWithKeyUnmarshal – error propagation when the
 * underlying monCommand fails (marshal error above)
 * ---------------------------------------------------------- */
func TestMonCommandWithKeyUnmarshal_PropagatesError(t *testing.T) {
	var bad unmarshalable
	var dst any
	err := monCommandWithKeyUnmarshal(nil, "someKey", bad, &dst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

//[TODO] nil cause
/*
func TestDumpMon_ErrorPropagation(t *testing.T) {
	_, err := dumpMon(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

func TestStatMon_ErrorPropagation(t *testing.T) {
	_, err := statMon(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}

func TestDumpOSD_ErrorPropagation(t *testing.T) {
	_, err := dumpOSD(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json")
}
*/

/* ---------------------------------------------------------- *
 * CephTime – JSON unmarshaller
 * ---------------------------------------------------------- */
func TestCephTime_UnmarshalJSON_EmptyAndNull(t *testing.T) {
	var ct CephTime
	// Empty string
	err := ct.UnmarshalJSON([]byte(`""`))
	assert.NoError(t, err)
	// Null
	err = ct.UnmarshalJSON([]byte(`null`))
	assert.NoError(t, err)
}

func TestCephTime_UnmarshalJSON_Valid(t *testing.T) {
	// Example timestamp used by the Ceph JSON dump:
	// "2023-01-02T15:04:05.000000-0700"
	const ts = `"2023-01-02T15:04:05.000000-0700"`
	var ct CephTime
	err := ct.UnmarshalJSON([]byte(ts))
	assert.NoError(t, err)

	expected, _ := time.Parse(`2006-01-02T15:04:05.000000-0700`, "2023-01-02T15:04:05.000000-0700")
	assert.Equal(t, expected, ct.Time)
}

/* ---------------------------------------------------------- *
 * CephSubvolumeTime – JSON unmarshaller
 * ---------------------------------------------------------- */
func TestCephSubvolumeTime_UnmarshalJSON_EmptyAndNull(t *testing.T) {
	var ct CephSubvolumeTime
	err := ct.UnmarshalJSON([]byte(`""`))
	assert.NoError(t, err)

	err = ct.UnmarshalJSON([]byte(`null`))
	assert.NoError(t, err)
}

func TestCephSubvolumeTime_UnmarshalJSON_Valid(t *testing.T) {
	// Example timestamp used by the FS dump:
	// "2023-01-02 15:04:05"
	const ts = `"2023-01-02 15:04:05"`
	var ct CephSubvolumeTime
	err := ct.UnmarshalJSON([]byte(ts))
	assert.NoError(t, err)

	expected, _ := time.Parse(`2006-01-02 15:04:05`, "2023-01-02 15:04:05")
	assert.Equal(t, expected, ct.Time)
}
