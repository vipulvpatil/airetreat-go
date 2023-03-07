package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_GetPlayerId(t *testing.T) {
	tests := []struct {
		name              string
		output            *pb.GetPlayerIdResponse
		playerCreatorMock storage.PlayerCreator
		errorExpected     bool
		errorString       string
	}{
		{
			name:              "test runs successfully",
			output:            &pb.GetPlayerIdResponse{PlayerId: "player_id1"},
			playerCreatorMock: &storage.PlayerCreatorMockSuccess{PlayerId: "player_id1"},
			errorExpected:     false,
			errorString:       "",
		},
		{
			name:              "errors if player creation fails",
			output:            nil,
			playerCreatorMock: &storage.PlayerCreatorMockFailure{},
			errorExpected:     true,
			errorString:       "unable to create player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithPlayerCreatorMock(
						tt.playerCreatorMock,
					),
				),
			})

			response, err := server.GetPlayerId(
				context.Background(),
				&pb.GetPlayerIdRequest{},
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
