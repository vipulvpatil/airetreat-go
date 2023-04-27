package server

import (
	"context"

	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	"google.golang.org/grpc"
)

// The calls to this service are authenticated using mutual TLS.
// This following interceptor adds a valid user if one exists
// on whose behalf the current request has been made.
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

func (s *AiRetreatGoService) PlayerIdValidatingInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	requestWithPlayerId, ok := req.(RequestWithPlayerId)
	if ok {
		playerId := requestWithPlayerId.GetPlayerId()
		if utilities.IsBlank(playerId) {
			return handler(ctx, req)
		}
		user, err := getUserFromContext(ctx)
		if err != nil {
			if utilities.ErrorIsUnauthenticated(err) && s.config.AllowUnauthed {
				player, err := s.storage.GetPlayer(playerId)
				if err != nil {
					return nil, &utilities.ResetPlayerError{}
				}
				if s.playerHasUser(player) {
					return nil, &utilities.ResetPlayerError{}
				}
			} else {
				return nil, err
			}
		} else {
			if !s.userPlayerIsNilOrSameAsPlayerId(user, playerId) {
				return nil, &utilities.ResetPlayerError{}
			}
		}
	}
	return handler(ctx, req)
}

type RequestWithPlayerId interface {
	GetPlayerId() string
}

func (s *AiRetreatGoService) userPlayerIsNilOrSameAsPlayerId(user *model.User, playerId string) bool {
	player, err := s.storage.GetPlayerForUserOrNil(user.GetId())
	if err != nil {
		return false
	}

	if player == nil {
		return true
	}

	return player.Id() == playerId
}

func (s *AiRetreatGoService) playerHasUser(player *model.Player) bool {
	return player.UserId() != nil
}
