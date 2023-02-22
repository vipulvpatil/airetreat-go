package storage

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func getGame(customDb customDbHandler, gameId string) (*model.Game, error) {
	if utilities.IsBlank(gameId) {
		return nil, errors.New("cannot getGame for a blank gameId")
	}

	var (
		opts           model.GameOptions
		stateHandledAt sql.NullTime
	)
	rows, err := customDb.Query(
		`SELECT
		g.id, g.state, g.current_turn_index, g.turn_order,
		g.state_handled, g.state_handled_at, g.state_total_time,
		g.created_at, g.updated_at,
		b.id, b.name, b.type, b.player_id, m.text
		FROM public."games" AS g
		LEFT JOIN public."bots" AS b ON b.game_id = g.id
		LEFT JOIN public."messages" AS m ON m.bot_id = b.id
		WHERE g.id = $1
		ORDER BY b.created_at ASC, b.id, m.created_at, m.id`, gameId,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select game")
	}
	defer rows.Close()

	botOptsMap := map[string]model.BotOptions{}
	botOptsOrderedIds := []string{}

	for rows.Next() {
		var botOpts model.BotOptions
		var playerId sql.NullString
		var message sql.NullString
		err := rows.Scan(
			&opts.Id,
			&opts.State,
			&opts.CurrentTurnIndex,
			pq.Array(&opts.TurnOrder),
			&opts.StateHandled,
			&stateHandledAt,
			&opts.StateTotalTime,
			&opts.CreatedAt,
			&opts.UpdatedAt,
			&botOpts.Id,
			&botOpts.Name,
			&botOpts.TypeOfBot,
			&playerId,
			&message,
		)

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
			if message.Valid {
				botOpts = botOptsMap[botOpts.Id]
				botOpts.Messages = append(botOpts.Messages, message.String)
				botOptsMap[botOpts.Id] = botOpts
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

func connectPlayerToBot(customDb customDbHandler, playerId, botId string) error {
	if utilities.IsBlank(playerId) {
		return errors.New("playerId cannot be blank")
	}

	if utilities.IsBlank(botId) {
		return errors.New("botId cannot be blank")
	}

	result, err := customDb.Exec(
		`UPDATE public."bots" SET "player_id" = $1, "type" = 'HUMAN' WHERE id = $2`, playerId, botId,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while connecting player to bot: %s %s", playerId, botId))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row after connecting player to bot: %s %s", playerId, botId))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError("No rows were affected when player was connected to Bot. This is highly unexpected.")
	}

	return nil
}

func updateGameStateToPlayersJoined(customDb customDbHandler, gameId string) error {
	if utilities.IsBlank(gameId) {
		return errors.New("gameId cannot be blank")
	}

	result, err := customDb.Exec(
		`UPDATE public."games"
		SET state = 'PLAYERS_JOINED'
		FROM public."games" AS g
		JOIN (
			SELECT count(id) AS human_bot_count, game_id
			FROM public."bots"
			WHERE game_id = $1
			AND type = 'HUMAN'
			GROUP BY game_id
			) AS b
		ON g.id = b.game_id
		WHERE g.id = $1
		AND
		b.human_bot_count = 2`, gameId,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while updating game state: %s", gameId))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row after updating game state: %s", gameId))
	}

	if rowsAffected > 1 {
		return utilities.NewBadError("More than one row was affected when game state was updated. This is highly unexpected.")
	}

	return nil
}
