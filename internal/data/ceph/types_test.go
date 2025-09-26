// types_test.go
package ceph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCephTime_UnmarshalJSON_AdditionalCases(t *testing.T) {
	var ct cephTime

	// Test with different time formats
	err := ct.UnmarshalJSON([]byte(`"2023-12-25T15:30:45.123456+0800"`))
	assert.NoError(t, err)
	assert.Equal(t, 2023, ct.Year())

	// Test with zero value
	err = ct.UnmarshalJSON([]byte(`"0000-01-01T00:00:00.000000+0000"`))
	assert.NoError(t, err)
}

func TestCephSubvolumeTime_UnmarshalJSON_AdditionalCases(t *testing.T) {
	var ct cephSubvolumeTime

	// Test with different time format
	err := ct.UnmarshalJSON([]byte(`"2023-12-25 15:30:45"`))
	assert.NoError(t, err)
	assert.Equal(t, 2023, ct.Year())

	// Test with midnight
	err = ct.UnmarshalJSON([]byte(`"2023-01-01 00:00:00"`))
	assert.NoError(t, err)
	assert.Equal(t, 2023, ct.Year())

	// Test with invalid JSON
	err = ct.UnmarshalJSON([]byte(`invalid-json`))
	assert.Error(t, err)
}
