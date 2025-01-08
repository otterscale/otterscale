package data

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewData, NewUserRepo)

type Data struct{}

func NewData() (*Data, func(), error) {
	d := &Data{}
	return d, func() {}, nil
}
