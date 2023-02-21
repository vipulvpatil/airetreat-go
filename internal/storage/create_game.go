package storage

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) CreateGame() error {
	id := s.IdGenerator.Generate()

	botNames := model.RandomBotNames()
	botOptionsList := []model.BotOptions{}
	bots := []*model.Bot{}

	for _, name := range botNames {
		botOpts := model.BotOptions{
			Id:        s.IdGenerator.Generate(),
			Name:      name,
			TypeOfBot: "AI",
		}
		botOptionsList = append(botOptionsList, botOpts)
		bot, err := model.NewBot(botOpts)
		if err != nil {
			return utilities.WrapBadError(err, "failed to create bot")
		}
		bots = append(bots, bot)
	}

	gameOption := model.GameOptions{
		Id:               id,
		State:            "STARTED",
		CurrentTurnIndex: 0,
		TurnOrder:        []string{"b", "p1", "b", "p2"},
		StateHandled:     false,
		Bots:             bots,
	}

	_, err := model.NewGame(gameOption)
	if err != nil {
		return utilities.WrapBadError(err, "failed to create game")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return utilities.WrapBadError(err, "failed to start db transaction")
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO public."games" (
			"id", "state", "current_turn_index", "turn_order", "state_handled"
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
		`,
		gameOption.Id, gameOption.State, gameOption.CurrentTurnIndex, pq.Array(gameOption.TurnOrder), gameOption.StateHandled,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, "dbError while inserting game and changing db")
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting game in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	for _, botOpts := range botOptionsList {
		result, err := tx.Exec(
			`INSERT INTO public."bots" (
				"id", "name", "type", "game_id"
			)
			VALUES (
				$1, $2, $3, $4
			)
			`,
			botOpts.Id, botOpts.Name, botOpts.TypeOfBot, gameOption.Id,
		)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return utilities.WrapBadError(err, "dbError while inserting bot and changing db")
		}

		if rowsAffected != 1 {
			return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting bot in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
		}
	}

	err = tx.Commit()
	return err
}
