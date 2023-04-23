package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/metadata"
)

func Test_GetPlayerId(t *testing.T) {
	tests := []struct {
		name               string
		output             *pb.GetPlayerIdResponse
		playerAccessorMock storage.PlayerAccessor
		errorExpected      bool
		errorString        string
	}{
		{
			name:               "test runs successfully",
			output:             &pb.GetPlayerIdResponse{PlayerId: "player_id1"},
			playerAccessorMock: &storage.PlayerAccessorMockSuccess{PlayerId: "player_id1"},
			errorExpected:      false,
			errorString:        "",
		},
		{
			name:               "errors if player creation fails",
			output:             nil,
			playerAccessorMock: &storage.PlayerAccessorMockFailure{},
			errorExpected:      true,
			errorString:        "unable to create player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithPlayerAccessorMock(
						tt.playerAccessorMock,
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

func Test_RegisterPlayerId(t *testing.T) {
	requestingUserId := "user_id1"
	requestingUserEmail := "user_email1"
	tests := []struct {
		name  string
		input struct {
			playerId            string
			requestingUserId    *string
			requestingUserEmail *string
		}
		output             *pb.RegisterPlayerIdResponse
		txShouldCommit     bool
		transactionMock    *storage.DatabaseTransactionMock
		playerAccessorMock storage.PlayerAccessor
		errorExpected      bool
		errorString        string
	}{
		{
			name: "errors if unable to get user from context",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:             nil,
			txShouldCommit:     false,
			transactionMock:    nil,
			playerAccessorMock: nil,
			errorExpected:      true,
			errorString:        "rpc error: code = Unauthenticated desc = could not retrieve requesting_user_email from context",
		},
		{
			name: "successfully connects the user to the player and returns the playerId",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:  true,
			transactionMock: &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return model.NewPlayer(model.PlayerOptions{
						Id: "player_id1",
					})
				},
				UpdatePlayerWithUserIdUsingTransactionInternal: func(playerId, userId string, transaction storage.DatabaseTransaction) error {
					assert.Equal(t, "player_id1", playerId)
					assert.Equal(t, "user_id1", userId)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "error if unable to get a transaction",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:     false,
			transactionMock:    nil,
			playerAccessorMock: nil,
			errorExpected:      true,
			errorString:        "unable to begin a db transaction",
		},
		{
			name: "errors if unable to get player from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:  false,
			transactionMock: &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return nil, errors.New("unable to get player")
				},
			},
			errorExpected: true,
			errorString:   "unable to get player",
		},
		{
			name: "errors if player is already registered with a user",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:  false,
			transactionMock: &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					userId2 := "user_id2"
					return model.NewPlayer(model.PlayerOptions{
						Id:     "player_id1",
						UserId: &userId2,
					})
				},
			},
			errorExpected: true,
			errorString:   "player is already registered",
		},
		{
			name: "errors if unable to update player",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:  false,
			transactionMock: &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return model.NewPlayer(model.PlayerOptions{
						Id: "player_id1",
					})
				},
				UpdatePlayerWithUserIdUsingTransactionInternal: func(playerId, userId string, transaction storage.DatabaseTransaction) error {
					assert.Equal(t, "player_id1", playerId)
					assert.Equal(t, "user_id1", userId)
					return errors.New("unable to update player")
				},
			},
			errorExpected: true,
			errorString:   "unable to update player",
		},
		{
			name: "successfully connects the user to the player and returns the playerId",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output: &pb.RegisterPlayerIdResponse{
				ConfirmedPlayerId: "player_id1",
			},
			txShouldCommit:  true,
			transactionMock: &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return model.NewPlayer(model.PlayerOptions{
						Id: "player_id1",
					})
				},
				UpdatePlayerWithUserIdUsingTransactionInternal: func(playerId, userId string, transaction storage.DatabaseTransaction) error {
					assert.Equal(t, "player_id1", playerId)
					assert.Equal(t, "user_id1", userId)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithPlayerAccessorMock(
						tt.playerAccessorMock,
					),
					storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
						Transaction: tt.transactionMock,
					}),
				),
			})

			ctx := context.Background()
			md := metadata.New(map[string]string{})
			if tt.input.requestingUserId != nil && tt.input.requestingUserEmail != nil {
				md.Append(requestingUserIdCtxKey, *tt.input.requestingUserId)
				md.Append(requestingUserEmailCtxKey, *tt.input.requestingUserEmail)
			}
			ctxWithMd := metadata.NewIncomingContext(ctx, md)

			response, err := server.RegisterPlayerId(
				ctxWithMd,
				&pb.RegisterPlayerIdRequest{
					PlayerId: tt.input.playerId,
				},
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
			if tt.transactionMock != nil {
				if tt.txShouldCommit {
					assert.True(t, tt.transactionMock.Committed, "transaction should have committed")
				} else {
					assert.True(t, tt.transactionMock.Rolledback, "transaction should have rolledback")
					assert.False(t, tt.transactionMock.Committed, "transaction should not have committed")
				}
			}
		})
	}
}
