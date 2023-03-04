package storage

import (
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/model"
)

func (s *Storage) GetUnhandledGameIdsForState(gameStateString string) []string {
	gameState := model.GameState(gameStateString)
	if !gameState.Valid() {
		return nil
	}

	rows, err := s.db.Query(
		`SELECT id
		FROM public."games"
		WHERE state = $1
		AND state_handled = false
		ORDER BY created_at DESC, id DESC
		`, gameState.String(),
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	gameIds := []string{}

	for rows.Next() {
		var gameId string
		err := rows.Scan(
			&gameId,
		)

		if err != nil {
			fmt.Printf("THIS IS BAD but no error: failed while scanning rows: %s\n", err)
			return nil
		}
		fmt.Println(gameId)
		gameIds = append(gameIds, gameId)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("THIS IS BAD but no error: failed to correctly go through bot rows: %s\n", err)
		return nil
	}
	return gameIds
}
