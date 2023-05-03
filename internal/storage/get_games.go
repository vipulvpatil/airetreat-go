package storage

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) GetGames(playerId string) ([]string, error) {
	if utilities.IsBlank(playerId) {
		return nil, errors.New("cannot GetGames for a blank playerId")
	}

	rows, err := s.db.Query(
		`SELECT game_id
		FROM public."bots"
		WHERE player_id = $1
		ORDER BY created_at DESC, id DESC
		`, playerId,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select games")
	}
	defer rows.Close()

	gameIds := []string{}

	for rows.Next() {
		var gameId string
		err := rows.Scan(
			&gameId,
		)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		gameIds = append(gameIds, gameId)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}
	return gameIds, nil
}

func (s *Storage) GetOldGames(gameExpiryDuration time.Duration) ([]string, error) {
	if gameExpiryDuration > -5*time.Minute {
		return nil, errors.New("invalid game expiry duration. Max acceptable time is -5 minutes.")
	}

	rows, err := s.db.Query(
		`SELECT id
		FROM public."games"
		WHERE created_at < $1
		ORDER BY created_at DESC, id DESC
		`, time.Now().Add(gameExpiryDuration),
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select games")
	}
	defer rows.Close()

	gameIds := []string{}

	for rows.Next() {
		var gameId string
		err := rows.Scan(
			&gameId,
		)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		gameIds = append(gameIds, gameId)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}
	return gameIds, nil
}

func (s *Storage) GetPublicJoinableGames() ([]string, error) {
	recent := time.Now().Add(-30 * time.Minute)

	rows, err := s.db.Query(
		`SELECT g.id, count(b.id)
		FROM public."games" AS g
		INNER JOIN public."bots" AS b ON b.game_id = g.id
		WHERE g.created_at > $1
		AND g.public = true
		AND g.state = 'STARTED'
		AND b.type = 'HUMAN'
		GROUP BY g.id
		ORDER BY g.created_at DESC, g.id DESC`,
		recent,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select games")
	}
	defer rows.Close()

	gameIds := []string{}

	for rows.Next() {
		var gameId string
		var humanBotCount int
		err := rows.Scan(
			&gameId,
			&humanBotCount,
		)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		if humanBotCount == 1 {
			gameIds = append(gameIds, gameId)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}
	return gameIds, nil
}
