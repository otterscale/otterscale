package model

import (
	"fmt"
	"strings"
)

type SubnetIPAddressType int64

const (
	// SubnetIPAddressTypeAuto is a SubnetIPAddressType of type Auto.
	SubnetIPAddressTypeAuto SubnetIPAddressType = iota
	// SubnetIPAddressTypeSticky is a SubnetIPAddressType of type Sticky.
	SubnetIPAddressTypeSticky
	// SubnetIPAddressTypeUserReserved is a SubnetIPAddressType of type UserReserved.
	SubnetIPAddressTypeUserReserved SubnetIPAddressType = iota + 2
	// SubnetIPAddressTypeDHCP is a SubnetIPAddressType of type DHCP.
	SubnetIPAddressTypeDHCP
	// SubnetIPAddressTypeDiscovered is a SubnetIPAddressType of type Discovered.
	SubnetIPAddressTypeDiscovered
)

var ErrInvalidSubnetIPAddressType = fmt.Errorf("not a valid SubnetIPAddressType, try [%s]", strings.Join(_SubnetIPAddressTypeNames, ", "))

const _SubnetIPAddressTypeName = "AutoStickyUserReservedDHCPDiscovered"

var _SubnetIPAddressTypeNames = []string{
	_SubnetIPAddressTypeName[0:4],
	_SubnetIPAddressTypeName[4:10],
	_SubnetIPAddressTypeName[10:22],
	_SubnetIPAddressTypeName[22:26],
	_SubnetIPAddressTypeName[26:36],
}

// SubnetIPAddressTypeNames returns a list of possible string values of SubnetIPAddressType.
func SubnetIPAddressTypeNames() []string {
	tmp := make([]string, len(_SubnetIPAddressTypeNames))
	copy(tmp, _SubnetIPAddressTypeNames)
	return tmp
}

// Values returns a list of possible string values of SubnetIPAddressType for ent EnumValues interface.
func (x SubnetIPAddressType) Values() []string {
	tmp := make([]string, len(_SubnetIPAddressTypeNames))
	copy(tmp, _SubnetIPAddressTypeNames)
	return tmp
}

var _SubnetIPAddressTypeMap = map[SubnetIPAddressType]string{
	SubnetIPAddressTypeAuto:         _SubnetIPAddressTypeName[0:4],
	SubnetIPAddressTypeSticky:       _SubnetIPAddressTypeName[4:10],
	SubnetIPAddressTypeUserReserved: _SubnetIPAddressTypeName[10:22],
	SubnetIPAddressTypeDHCP:         _SubnetIPAddressTypeName[22:26],
	SubnetIPAddressTypeDiscovered:   _SubnetIPAddressTypeName[26:36],
}

// String implements the Stringer interface.
func (x SubnetIPAddressType) String() string {
	if str, ok := _SubnetIPAddressTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("SubnetIPAddressType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x SubnetIPAddressType) IsValid() bool {
	_, ok := _SubnetIPAddressTypeMap[x]
	return ok
}

var _SubnetIPAddressTypeValue = map[string]SubnetIPAddressType{
	_SubnetIPAddressTypeName[0:4]:   SubnetIPAddressTypeAuto,
	_SubnetIPAddressTypeName[4:10]:  SubnetIPAddressTypeSticky,
	_SubnetIPAddressTypeName[10:22]: SubnetIPAddressTypeUserReserved,
	_SubnetIPAddressTypeName[22:26]: SubnetIPAddressTypeDHCP,
	_SubnetIPAddressTypeName[26:36]: SubnetIPAddressTypeDiscovered,
}

// ParseSubnetIPAddressType attempts to convert a string to a SubnetIPAddressType.
func ParseSubnetIPAddressType(name string) (SubnetIPAddressType, error) {
	if x, ok := _SubnetIPAddressTypeValue[name]; ok {
		return x, nil
	}
	return SubnetIPAddressType(0), fmt.Errorf("%s is %w", name, ErrInvalidSubnetIPAddressType)
}

// MarshalText implements the text marshaller method.
func (x SubnetIPAddressType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *SubnetIPAddressType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseSubnetIPAddressType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
