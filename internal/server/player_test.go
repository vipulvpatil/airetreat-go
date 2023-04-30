package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/metadata"
)

func Test_SyncPlayerData(t *testing.T) {
	requestingUserId := "user_id1"
	requestingUserEmail := "user_email1"
	playerWithUser, _ := model.NewPlayer(model.PlayerOptions{Id: "player_id1", UserId: &requestingUserId})
	playerWithoutUser, _ := model.NewPlayer(model.PlayerOptions{Id: "player_id1", UserId: nil})
	anotherUserId := "user_id2"
	tests := []struct {
		name  string
		input struct {
			playerId            string
			requestingUserId    *string
			requestingUserEmail *string
		}
		output              *pb.SyncPlayerDataResponse
		allowUnauthedConfig bool
		txShouldCommit      bool
		transactionMock     *storage.DatabaseTransactionMock
		playerAccessorMock  storage.PlayerAccessor
		errorExpected       bool
		errorString         string
	}{
		{
			name: "errors if unable to get user from context and AllowUnauthed config is false",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  nil,
			errorExpected:       true,
			errorString:         "rpc error: code = Unauthenticated desc = could not retrieve requesting_user_email from context",
		},
		{
			name: "returns player if it exists for the user",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              &pb.SyncPlayerDataResponse{PlayerId: playerWithUser.Id(), Connected: true},
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockSuccess{PlayerId: "player_id1", UserId: &requestingUserEmail},
			errorExpected:       false,
			errorString:         "",
		},
		{
			name: "errors if unable to get the user from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockFailure{},
			errorExpected:       true,
			errorString:         "unable to get player",
		},
		{
			name: "when user does not have connected player and request playerId is not blank, errors if unable to get transaction",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
			},
			errorExpected: true,
			errorString:   "unable to begin a db transaction",
		},
		{
			name: "when user does not have connected player and request playerId is not blank, errors if unable to get player",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return nil, errors.New("unable to get player")
				},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: unknown error",
		},
		{
			name: "when user does not have connected player and request playerId is not blank, errors if player has a different connected user",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return model.NewPlayer(model.PlayerOptions{Id: "player_id1", UserId: &anotherUserId})
				},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: unknown error",
		},
		{
			name: "when user does not have connected player and request playerId is not blank, errors if connecting user and player fails",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              nil,
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return playerWithoutUser, nil
				},
				UpdatePlayerWithUserIdUsingTransactionInternal: func(playerId, userId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return nil, errors.New("unable to connect player to user")
				},
			},
			errorExpected: true,
			errorString:   "unable to connect player to user",
		},
		{
			name: "when user does not have connected player and request playerId is not blank, returns connecting player",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              &pb.SyncPlayerDataResponse{PlayerId: playerWithUser.Id(), Connected: true},
			allowUnauthedConfig: false,
			txShouldCommit:      true,
			transactionMock:     &storage.DatabaseTransactionMock{},
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
				GetPlayerUsingTransactionInternal: func(playerId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return playerWithoutUser, nil
				},
				UpdatePlayerWithUserIdUsingTransactionInternal: func(playerId, userId string, transaction storage.DatabaseTransaction) (*model.Player, error) {
					return playerWithUser, nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "when user does not have connected player and request playerId is blank, creates and returns a new connected player",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "",
				requestingUserId:    &requestingUserId,
				requestingUserEmail: &requestingUserEmail,
			},
			output:              &pb.SyncPlayerDataResponse{PlayerId: "new_player_id1", Connected: true},
			allowUnauthedConfig: false,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock: &storage.PlayerAccessorMockConfigurable{
				GetPlayerForUserOrNilInternal: func(userId string) (*model.Player, error) {
					return nil, nil
				},
				CreatePlayerForUserInternal: func(userId string) (*model.Player, error) {
					return model.NewPlayer(model.PlayerOptions{Id: "new_player_id1", UserId: &requestingUserId})
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "when no user in context and request playerId is non blank, errors if unable to get player from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              nil,
			allowUnauthedConfig: true,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockFailure{},
			errorExpected:       true,
			errorString:         "unable to get player",
		},
		{
			name: "when no user in context and request playerId is non blank, errors if player returned from storage has a User associated",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              nil,
			allowUnauthedConfig: true,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockSuccess{PlayerId: "player_id1", UserId: &anotherUserId},
			errorExpected:       true,
			errorString:         "reset player data",
		},
		{
			name: "when no user in context and request playerId is non blank, returns the player from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "player_id1",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              &pb.SyncPlayerDataResponse{PlayerId: "player_id1", Connected: false},
			allowUnauthedConfig: true,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockSuccess{PlayerId: "player_id1"},
			errorExpected:       false,
			errorString:         "",
		},
		{
			name: "when no user in context and request playerId is blank, errors if unable to create player from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              nil,
			allowUnauthedConfig: true,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockFailure{},
			errorExpected:       true,
			errorString:         "unable to create player",
		},
		{
			name: "when no user in context and request playerId is blank, returns the newly created player from storage",
			input: struct {
				playerId            string
				requestingUserId    *string
				requestingUserEmail *string
			}{
				playerId:            "",
				requestingUserId:    nil,
				requestingUserEmail: nil,
			},
			output:              &pb.SyncPlayerDataResponse{PlayerId: "player_id1", Connected: false},
			allowUnauthedConfig: true,
			txShouldCommit:      false,
			transactionMock:     nil,
			playerAccessorMock:  &storage.PlayerAccessorMockSuccess{PlayerId: "player_id1"},
			errorExpected:       false,
			errorString:         "",
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
				Config: &config.Config{
					AllowUnauthed: tt.allowUnauthedConfig,
				},
				Logger: &utilities.NullLogger{},
			})

			ctx := context.Background()
			md := metadata.New(map[string]string{})
			if tt.input.requestingUserId != nil && tt.input.requestingUserEmail != nil {
				md.Append(requestingUserIdCtxKey, *tt.input.requestingUserId)
				md.Append(requestingUserEmailCtxKey, *tt.input.requestingUserEmail)
			}
			ctxWithMd := metadata.NewIncomingContext(ctx, md)

			response, err := server.SyncPlayerData(
				ctxWithMd,
				&pb.SyncPlayerDataRequest{
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
