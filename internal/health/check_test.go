package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_Check(t *testing.T) {
	t.Run("test runs successfully", func(t *testing.T) {
		server := AiRetreatGoHealthService{}

		response, err := server.Check(
			context.Background(),
			&pb.CheckRequest{},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, response, &pb.CheckResponse{})
	})
}
