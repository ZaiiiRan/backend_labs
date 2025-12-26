package utils

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func GetGrpcErrStatus(err error) (*status.Status, error) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, fmt.Errorf("not a grpc error: %w", err)
	}
	return st, nil
}
