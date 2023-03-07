package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_SendMessage(t *testing.T) {
	tests := []struct {
		name               string
		input              *pb.SendMessageRequest
		output             *pb.SendMessageResponse
		messageCreatorMock storage.MessageCreator
		errorExpected      bool
		errorString        string
	}{
		{
			name: "test runs successfully",
			input: &pb.SendMessageRequest{
				BotId: "bot_id1",
				Text:  "message",
			},
			output:             &pb.SendMessageResponse{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
			errorExpected:      false,
			errorString:        "",
		},
		{
			name: "errors if message creation fails",
			input: &pb.SendMessageRequest{
				BotId: "bot_id1",
				Text:  "message",
			},
			output:             nil,
			messageCreatorMock: &storage.MessageCreatorMockFailure{},
			errorExpected:      true,
			errorString:        "unable to create message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithMessageCreatorMock(
						tt.messageCreatorMock,
					),
				),
			})

			response, err := server.SendMessage(
				context.Background(),
				tt.input,
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
