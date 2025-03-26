package model

import (
	"fmt"
	"strings"
)

type AllocType int64

const (
	// AllocTypeAutomatic is a AllocType of type Automatic.
	AllocTypeAutomatic AllocType = iota
	// AllocTypeSticky is a AllocType of type Sticky.
	AllocTypeSticky
	// AllocTypeUserReserved is a AllocType of type User Reserved.
	AllocTypeUserReserved AllocType = iota + 2
	// AllocTypeDHCP is a AllocType of type DHCP.
	AllocTypeDHCP
	// AllocTypeDiscovered is a AllocType of type Discovered.
	AllocTypeDiscovered
)

var ErrInvalidAllocType = fmt.Errorf("not a valid AllocType, try [%s]", strings.Join(_AllocTypeNames, ", "))

const _AllocTypeName = "AutomaticStickyUser ReservedDHCPDiscovered"

var _AllocTypeNames = []string{
	_AllocTypeName[0:9],
	_AllocTypeName[9:15],
	_AllocTypeName[15:28],
	_AllocTypeName[28:32],
	_AllocTypeName[32:42],
}

// AllocTypeNames returns a list of possible string values of AllocType.
func AllocTypeNames() []string {
	tmp := make([]string, len(_AllocTypeNames))
	copy(tmp, _AllocTypeNames)
	return tmp
}

// Values returns a list of possible string values of AllocType for ent EnumValues interface.
func (x AllocType) Values() []string {
	tmp := make([]string, len(_AllocTypeNames))
	copy(tmp, _AllocTypeNames)
	return tmp
}

var _AllocTypeMap = map[AllocType]string{
	AllocTypeAutomatic:    _AllocTypeName[0:9],
	AllocTypeSticky:       _AllocTypeName[9:15],
	AllocTypeUserReserved: _AllocTypeName[15:28],
	AllocTypeDHCP:         _AllocTypeName[28:32],
	AllocTypeDiscovered:   _AllocTypeName[32:42],
}

// String implements the Stringer interface.
func (x AllocType) String() string {
	if str, ok := _AllocTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("AllocType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x AllocType) IsValid() bool {
	_, ok := _AllocTypeMap[x]
	return ok
}

var _AllocTypeValue = map[string]AllocType{
	_AllocTypeName[0:9]:   AllocTypeAutomatic,
	_AllocTypeName[9:15]:  AllocTypeSticky,
	_AllocTypeName[15:28]: AllocTypeUserReserved,
	_AllocTypeName[28:32]: AllocTypeDHCP,
	_AllocTypeName[32:42]: AllocTypeDiscovered,
}

// ParseAllocType attempts to convert a string to a AllocType.
func ParseAllocType(name string) (AllocType, error) {
	if x, ok := _AllocTypeValue[name]; ok {
		return x, nil
	}
	return AllocType(0), fmt.Errorf("%s is %w", name, ErrInvalidAllocType)
}

// MarshalText implements the text marshaller method.
func (x AllocType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *AllocType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseAllocType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
