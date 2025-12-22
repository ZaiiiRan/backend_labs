package utils

import "google.golang.org/grpc/status"

func IsGrpcError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := status.FromError(err)
	return ok
}
