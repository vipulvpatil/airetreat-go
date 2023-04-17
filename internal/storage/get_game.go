package storage

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) GetGame(gameId string) (*model.Game, error) {
	return getGameUsingCustomDbHandler(s.db, gameId, false)
}

func (s *Storage) GetGameUsingTransaction(gameId string, transaction DatabaseTransaction) (*model.Game, error) {
	return getGameUsingCustomDbHandler(transaction, gameId, true)
}

func getGameUsingCustomDbHandler(customDb customDbHandler, gameId string, exclusiveLock bool) (*model.Game, error) {
	if utilities.IsBlank(gameId) {
		return nil, errors.New("cannot getGame for a blank gameId")
	}

	var (
		opts           model.GameOptions
		stateHandledAt sql.NullTime
	)

	queryWithoutLock := `SELECT
	g.id, g.state, g.current_turn_index, g.turn_order,
	g.state_handled, g.state_handled_at, g.state_total_time,
	g.last_question, g.last_question_target_bot_id,
	g.result, g.winning_bot_id,
	g.created_at, g.updated_at,
	b.id, b.name, b.type, b.player_id, b.help_count,
	m.source_bot_id, m.target_bot_id, m.text, m.created_at, m.type
	FROM public."games" AS g
	LEFT JOIN public."bots" AS b ON b.game_id = g.id
	LEFT JOIN public."messages" AS m ON m.target_bot_id = b.id
	WHERE g.id = $1
	ORDER BY b.created_at ASC, b.id ASC, m.created_at ASC, m.id ASC`

	queryWithLock := `SELECT
	g.id, g.state, g.current_turn_index, g.turn_order,
	g.state_handled, g.state_handled_at, g.state_total_time,
	g.last_question, g.last_question_target_bot_id,
	g.result, g.winning_bot_id,
	g.created_at, g.updated_at,
	b.id, b.name, b.type, b.player_id, b.help_count,
	m.source_bot_id, m.target_bot_id, m.text, m.created_at, m.type
	FROM public."games" AS g
	LEFT JOIN public."bots" AS b ON b.game_id = g.id
	LEFT JOIN public."messages" AS m ON m.target_bot_id = b.id
	WHERE g.id = $1
	ORDER BY b.created_at ASC, b.id ASC, m.created_at ASC, m.id ASC
	FOR UPDATE OF g`

	var query string

	if exclusiveLock {
		query = queryWithLock
	} else {
		query = queryWithoutLock
	}

	rows, err := customDb.Query(query, gameId)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select game")
	}
	defer rows.Close()

	botOptsMap := map[string]model.BotOptions{}
	botOptsOrderedIds := []string{}

	for rows.Next() {
		var botOpts model.BotOptions
		var playerId sql.NullString
		var messageSourceBotId sql.NullString
		var messageTargetBotId sql.NullString
		var messageText sql.NullString
		var messageType sql.NullString
		var lastQuestion sql.NullString
		var lastQuestionTargetBotId sql.NullString
		var result sql.NullString
		var winningBotId sql.NullString
		var messageCreatedAt sql.NullTime
		err := rows.Scan(
			&opts.Id,
			&opts.State,
			&opts.CurrentTurnIndex,
			pq.Array(&opts.TurnOrder),
			&opts.StateHandled,
			&stateHandledAt,
			&opts.StateTotalTime,
			&lastQuestion,
			&lastQuestionTargetBotId,
			&result,
			&winningBotId,
			&opts.CreatedAt,
			&opts.UpdatedAt,
			&botOpts.Id,
			&botOpts.Name,
			&botOpts.TypeOfBot,
			&playerId,
			&botOpts.HelpCount,
			&messageSourceBotId,
			&messageTargetBotId,
			&messageText,
			&messageCreatedAt,
			&messageType,
		)

		if lastQuestion.Valid {
			opts.LastQuestion = lastQuestion.String
		}
		if lastQuestionTargetBotId.Valid {
			opts.LastQuestionTargetBotId = lastQuestionTargetBotId.String
		}
		if result.Valid {
			opts.Result = result.String
		}
		if winningBotId.Valid {
			opts.WinningBotId = winningBotId.String
		}

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		if !utilities.IsBlank(botOpts.Id) {
			if playerId.Valid {
				player, err := model.NewPlayer(model.PlayerOptions{Id: playerId.String})
				if err != nil {
					return nil, utilities.WrapBadError(err, "failed to create player")
				}
				botOpts.ConnectedPlayer = player
			}
			_, ok := botOptsMap[botOpts.Id]
			if !ok {
				botOptsOrderedIds = append(botOptsOrderedIds, botOpts.Id)
				botOptsMap[botOpts.Id] = botOpts
			}
			if messageText.Valid {
				message := model.Message{
					SourceBotId: messageSourceBotId.String,
					TargetBotId: messageTargetBotId.String,
					Text:        messageText.String,
					CreatedAt:   messageCreatedAt.Time,
					MessageType: messageType.String,
				}
				opts.Messages = append(opts.Messages, &message)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}

	for _, botOptsId := range botOptsOrderedIds {
		bot, err := model.NewBot(botOptsMap[botOptsId])
		if err != nil {
			return nil, utilities.WrapBadError(err, "failed to create bot")
		}
		opts.Bots = append(opts.Bots, bot)
	}

	if stateHandledAt.Valid {
		opts.StateHandledAt = &stateHandledAt.Time
	}

	if utilities.IsBlank(opts.Id) {
		return nil, errors.Errorf("game not found: %s", gameId)
	}

	game, err := model.NewGame(opts)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to create game")
	}
	return game, nil
}
