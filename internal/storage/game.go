package storage

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type GameAccessor interface {
	CreateGame() error
	JoinGame(gameId string, player *model.Player) error
}

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
		b.id, b.name, b.type, b.player_id
		FROM public."games" AS g
		LEFT JOIN public."bots" AS b ON b.game_id = g.id
		WHERE g.id = $1`, gameId,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select game")
	}
	defer rows.Close()

	for rows.Next() {
		var botOpts model.BotOptions
		var playerId sql.NullString
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
		)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		if !utilities.IsBlank(botOpts.Id) {
			bot, err := model.NewBot(botOpts)
			if err != nil {
				return nil, utilities.WrapBadError(err, "failed to create bot")
			}
			if playerId.Valid {
				player, err := model.NewPlayer(model.PlayerOptions{Id: playerId.String})
				if err != nil {
					return nil, utilities.WrapBadError(err, "failed to create player")
				}
				err = bot.ConnectPlayer(player)
				if err != nil {
					return nil, utilities.WrapBadError(err, "failed to connect player")
				}
			}
			opts.Bots = append(opts.Bots, bot)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}

	if stateHandledAt.Valid {
		opts.StateHandledAt = &stateHandledAt.Time
	}

	if utilities.IsBlank(opts.Id) {
		return nil, errors.Errorf("game not found: %s", gameId)
	}

	game, err := model.NewGame(opts)
	fmt.Println(err)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to create game")
	}
	return game, nil
}

func getRandomBot(bots []*model.Bot) (*model.Bot, error) {
	if len(bots) == 0 {
		return nil, errors.Errorf("Cannot get random bot from an empty list")
	}

	rand.Shuffle(len(bots), func(i, j int) {
		bots[i], bots[j] = bots[j], bots[i]
	})

	return bots[0], nil
}

func getAllAiBots(customDb customDbHandler, gameId string) ([]*model.Bot, error) {
	if utilities.IsBlank(gameId) {
		return nil, errors.New("gameId cannot be blank")
	}

	rows, err := customDb.Query(
		`SELECT id, name, type, player_id FROM public."bots" WHERE game_id = $1`, gameId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("GetAIBots %s: no bots for game", gameId)
		}
		return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError %s", gameId))
	}
	defer rows.Close()

	bots := []*model.Bot{}
	for rows.Next() {
		var (
			botOptions model.BotOptions
			player     sql.NullString
		)
		err := rows.Scan(&botOptions.Id, &botOptions.Name, &botOptions.TypeOfBot, player)
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while scanning row of bots for game: %s", gameId))
		}
		if player.Valid {
			return nil, utilities.NewBadError(fmt.Sprintf("AI bot has connected player for bot: %v", botOptions.Id))
		}

		bot, err := model.NewBot(botOptions)
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("bot returned by db was inconsistent :%s", botOptions.Id))
		}
		bots = append(bots, bot)
	}

	return bots, nil
}

func connectPlayerToBot(customDb customDbHandler, playerId, botId string) error {
	if utilities.IsBlank(playerId) {
		return errors.New("playerId cannot be blank")
	}

	if utilities.IsBlank(botId) {
		return errors.New("botId cannot be blank")
	}

	result, err := customDb.Exec(
		`UPDATE public."bots" (player_id) VALUES ($1) WHERE id = $2`, playerId, botId,
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
