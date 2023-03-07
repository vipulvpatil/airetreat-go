package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/metadata"
)

func contextWithPrefilledRequestingUser() context.Context {
	return metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			requestingUserIdCtxKey, "internalUserId1",
			requestingUserEmailCtxKey, "test@example.com",
		),
	)
}

func Test_Test(t *testing.T) {
	t.Run("test runs successfully", func(t *testing.T) {
		server, _ := NewServer(ServerDependencies{})

		response, err := server.Test(
			contextWithPrefilledRequestingUser(),
			&pb.TestRequest{Test: "test_string"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, response, &pb.TestResponse{Test: "success: test_string"})
	})
}
