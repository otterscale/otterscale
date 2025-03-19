package app

import (
	"encoding/base64"
	"strconv"

	"github.com/google/wire"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ProviderSet = wire.NewSet(NewStackApp)

var (
	defaultPageSize     = 200
	maxPageSize         = 10000
	ErrInvalidPageSize  = status.Errorf(codes.InvalidArgument, "page size cannot be less than zero")
	ErrInvalidPageToken = status.Errorf(codes.InvalidArgument, "page token is invalid")
	ErrInvalidQuery     = status.Errorf(codes.InvalidArgument, "query is invalid")
)

func getPageSize(pageSize int32) (int, error) {
	switch {
	case pageSize < 0:
		return -1, ErrInvalidPageSize
	case pageSize == 0:
		return defaultPageSize, nil
	case int(pageSize) > maxPageSize:
		return maxPageSize, nil
	}
	return int(pageSize), nil
}

func getPageToken(pageToken string) (int, error) {
	if pageToken == "" {
		return -1, nil
	}
	bs, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return -1, ErrInvalidPageToken
	}
	return strconv.Atoi(string(bs))
}
