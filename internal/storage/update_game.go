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

func sqlAndArgsForUpdate(updateOpts GameUpdateOptions) (string, []interface{}) {
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

	finalSql := strings.Join(setSqls, ", ")
	return finalSql, args
}

func (s *Storage) UpdateGameState(gameId string, updateOpts GameUpdateOptions) error {
	updateSqlSetPart, args := sqlAndArgsForUpdate(updateOpts)
	if len(args) == 0 {
		return errors.New("no update options provided")
	}

	argsWithGameId := append(args, gameId)

	updateSql := fmt.Sprintf("UPDATE public.\"games\" SET %s WHERE \"id\" = $%d", updateSqlSetPart, len(args)+1)
	result, err := s.db.Exec(updateSql, argsWithGameId...)

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
