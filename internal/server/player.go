package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

// TODO: this function is overly complicated and needs to be simplified.
func (s *AiRetreatGoService) SyncPlayerData(ctx context.Context, req *pb.SyncPlayerDataRequest) (*pb.SyncPlayerDataResponse, error) {
	user, err := s.getUserFromContextIfPresent(ctx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	playerId := req.GetPlayerId()
	var player *model.Player
	if user != nil {
		player, err = s.getNewOrExistingPlayerForUser(user.GetId(), playerId)
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}
	} else if !utilities.IsBlank(playerId) {
		player, err = s.storage.GetPlayer(playerId)
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}

		if player.UserId() != nil {
			return nil, &utilities.ResetPlayerError{}
		}
	} else {
		player, err = s.storage.CreatePlayer()
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}
	}

	response := pb.SyncPlayerDataResponse{
		PlayerId:  player.Id(),
		Connected: player.UserId() != nil,
	}

	return &response, err
}

func (s *AiRetreatGoService) getUserFromContextIfPresent(ctx context.Context) (*model.User, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		if utilities.ErrorIsUnauthenticated(err) && s.config.AllowUnauthed {
			return nil, nil
		}
		s.logger.LogError(err)
		return nil, err
	}
	return user, nil
}

func (s *AiRetreatGoService) getNewOrExistingPlayerForUser(userId string, playerId string) (*model.Player, error) {
	if utilities.IsBlank(userId) {
		err := errors.New("userId cannot be blank")
		s.logger.LogError(err)
		return nil, err
	}

	player, err := s.storage.GetPlayerForUserOrNil(userId)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	if player != nil {
		// TODO: Not sure if this will be a problem. Currently we simply ignore the playedId in the request, if the user has a player connected already
		return player, nil
	}

	if !utilities.IsBlank(playerId) {
		tx, err := s.storage.BeginTransaction()
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}
		defer tx.Rollback()
		player, err = s.storage.GetPlayerUsingTransaction(playerId, tx)
		// TODO: Rethink this. This can be used to find playerIds that are connected to some user in our system. Not sure if that is a security risk. Sending unknown error for now.
		if err != nil {
			s.logger.LogError(err)
			return nil, utilities.NewBadError("unknown error")
		}

		// TODO: Rethink this. This can be used to find playerIds that are connected to some user in our system. Not sure if that is a security risk. Sending unknown error for now.
		if player.UserId() != nil {
			s.logger.LogError(utilities.NewBadError("unknown error"))
			return nil, utilities.NewBadError("unknown error")
		}

		player, err = s.storage.UpdatePlayerWithUserIdUsingTransaction(player.Id(), userId, tx)
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			s.logger.LogError(err)
			return nil, err
		}

		return player, nil
	} else {
		return s.storage.CreatePlayerForUser(userId)
	}
}
