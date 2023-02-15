package storage

import "github.com/vipulvpatil/airetreat-go/internal/model"

type GameAccessor interface {
	CreateGame() model.Game
}

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
			return err
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
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO public."games" (
			"id", "state", "current_turn_index", "turn_order", "state_handled"
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
		`,
		gameOption.Id, gameOption.State, gameOption.CurrentTurnIndex, gameOption.TurnOrder, gameOption.StateHandled,
	)
	if err != nil {
		return err
	}

	for _, botOpts := range botOptionsList {
		_, err = tx.Exec(
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
	}

	err = tx.Commit()

	return err
}
