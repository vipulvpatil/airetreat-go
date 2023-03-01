package storage

import (
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
