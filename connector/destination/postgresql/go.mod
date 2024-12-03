module github.com/openhdc/openhdc-connector

go 1.23.3

replace github.com/openhdc/openhdc v0.0.0-20241129065929-69142d07d680 => ../../../../openhdc

require (
	github.com/google/wire v0.6.0
	github.com/openhdc/openhdc v0.0.0-20241129065929-69142d07d680
	go.uber.org/automaxprocs v1.6.0
)

require (
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.68.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
