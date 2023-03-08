package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type GameUpdateOptions struct {
	State                   *string
	CurrentTurnIndex        *int64
	TurnOrder               []string
	StateHandled            *bool
	StateHandledAt          *time.Time
	LastQuestion            *string
	LastQuestionTargetBotId *string
	StateTotalTime          *int64
}

func (s *Storage) UpdateGameState(gameId string, updateOpts GameUpdateOptions) error {
	return updateGameState(s.db, gameId, updateOpts)
}

func (s *Storage) UpdateGameStateUsingTransaction(gameId string, updateOpts GameUpdateOptions, transaction DatabaseTransaction) error {
	return updateGameState(transaction, gameId, updateOpts)
}

func updateGameState(customDb customDbHandler, gameId string, updateOpts GameUpdateOptions) error {
	updateSqlsPart, args := sqlAndArgsForUpdate(updateOpts)
	if len(args) == 0 {
		return errors.New("no update options provided")
	}
	updateSqlsPart, args = autoAddUpdateTimeStamp(updateSqlsPart, args)
	updateSqlSetPart := strings.Join(updateSqlsPart, ", ")
	argsWithGameId := append(args, gameId)
	updateSql := fmt.Sprintf("UPDATE public.\"games\" SET %s WHERE \"id\" = $%d", updateSqlSetPart, len(args)+1)
	result, err := customDb.Exec(updateSql, argsWithGameId...)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, "dbError while updating game and changing db")
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when updating game in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}
	return nil
}

func (s *Storage) UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(gameId string, transaction DatabaseTransaction) error {
	return updateGameStateIfEnoughPlayersHaveJoined(transaction, gameId)
}

func sqlAndArgsForUpdate(updateOpts GameUpdateOptions) ([]string, []interface{}) {
	// Changes to this function may result in SQL injection issues. Be careful while modifying.
	setSqls := []string{}
	args := []interface{}{}
	index := 1

	if updateOpts.State != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"state\" = $%d", index))
		args = append(args, *updateOpts.State)
		index++
	}
	if updateOpts.CurrentTurnIndex != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"current_turn_index\" = $%d", index))
		args = append(args, *updateOpts.CurrentTurnIndex)
		index++
	}
	if updateOpts.TurnOrder != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"turn_order\" = $%d", index))
		args = append(args, pq.Array(updateOpts.TurnOrder))
		index++
	}
	if updateOpts.StateHandled != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"state_handled\" = $%d", index))
		args = append(args, *updateOpts.StateHandled)
		index++
	}
	if updateOpts.StateHandledAt != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"state_handled_at\" = $%d", index))
		args = append(args, *updateOpts.StateHandledAt)
		index++
	}
	if updateOpts.LastQuestion != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"last_question\" = $%d", index))
		args = append(args, *updateOpts.LastQuestion)
		index++
	}
	if updateOpts.LastQuestionTargetBotId != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"last_question_target_bot_id\" = $%d", index))
		args = append(args, *updateOpts.LastQuestionTargetBotId)
		index++
	}
	if updateOpts.StateTotalTime != nil {
		setSqls = append(setSqls, fmt.Sprintf("\"state_total_time\" = $%d", index))
		args = append(args, *updateOpts.StateTotalTime)
		index++
	}

	return setSqls, args
}

func autoAddUpdateTimeStamp(setSqls []string, args []interface{}) ([]string, []interface{}) {
	if len(args) > 0 {
		setSqls = append(setSqls, fmt.Sprintf("\"updated_at\" = $%d", len(args)+1))
		args = append(args, time.Now())
	}
	return setSqls, args
}

func updateGameStateIfEnoughPlayersHaveJoined(customDb customDbHandler, gameId string) error {
	if utilities.IsBlank(gameId) {
		return errors.New("gameId cannot be blank")
	}

	result, err := customDb.Exec(
		`WITH selected_games AS (
			SELECT g.id, count(b.id) AS human_bot_count
			FROM public."games" AS g
			JOIN public."bots" AS b ON g.id = b.game_id
			WHERE g.id = $1
			AND b.type = 'HUMAN'
			GROUP BY g.id
		)
		UPDATE public."games" AS games
		SET state = 'PLAYERS_JOINED', updated_at = $2
		FROM selected_games
		WHERE games.id = selected_games.id
		AND human_bot_count = 2`,
		gameId, time.Now(),
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while updating game state: %s", gameId))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row after updating game state: %s", gameId))
	}

	if rowsAffected > 1 {
		return utilities.NewBadError(fmt.Sprintf("More than one row (%d) was affected when game state was updated. This is highly unexpected.", rowsAffected))
	}

	return nil
}
