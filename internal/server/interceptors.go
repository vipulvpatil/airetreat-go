package server

import (
	"context"

	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	"google.golang.org/grpc"
)

// The calls to this service are authenticated using mutual TLS.
// This following interceptor ensures there is a valid user on whose behalf
// the current request has been made.
func (s *AiRetreatGoService) RequestingUserInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	updatedCtx, err := contextWithUserData(ctx, s.storage)
	if err != nil {
		if utilities.ErrorIsUnauthenticated(err) && s.config.AllowUnauthed {
			return handler(ctx, req)
		} else {
			return nil, err
		}
	} else {
		return handler(updatedCtx, req)
	}
}
