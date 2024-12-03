package workload

import (
	"errors"
	"fmt"
)

type Kind int

const (
	KindSource Kind = iota
	KindDestination
	KindTransformer
)

var ErrInvalidKind = errors.New("invalid kind")

var KindMap = map[Kind]string{
	KindSource:      "source",
	KindDestination: "destination",
	KindTransformer: "transformer",
}

func (k Kind) String() string {
	if val, ok := KindMap[k]; ok {
		return val
	}
	return fmt.Sprintf("Kind(%d)", k)
}

func (k Kind) MarshalYAML() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k *Kind) UnmarshalYAML(str []byte) error {
	tmp, err := ParseKind(string(str))
	if err != nil {
		return err
	}
	*k = tmp
	return nil
}

func ParseKind(str string) (Kind, error) {
	for key, val := range KindMap {
		if val == str {
			return Kind(key), nil
		}
	}
	return Kind(0), fmt.Errorf("invalid kind %s", str)
}
