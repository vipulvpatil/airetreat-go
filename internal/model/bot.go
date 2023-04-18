package model

import (
	"math/rand"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

const totalNumberOfBotsPerGame = 5

type Bot struct {
	id        string
	name      string
	typeOfBot botType
	player    *Player
	helpCount int64
}

type BotOptions struct {
	Id              string
	Name            string
	TypeOfBot       string
	ConnectedPlayer *Player
	HelpCount       int64
}

func NewBot(opts BotOptions) (*Bot, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create bot with an empty id")
	}
	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create bot with an empty name")
	}

	defaultBotType := ai
	if utilities.IsBlank(opts.TypeOfBot) {
		opts.TypeOfBot = defaultBotType.String()
	}
	typeOfBot := BotType(opts.TypeOfBot)
	if !typeOfBot.Valid() {
		return nil, errors.New("cannot create bot with an invalid botType")
	}

	if opts.ConnectedPlayer != nil && typeOfBot != human {
		return nil, errors.New("cannot create a bot of non-human type with a connected Player")
	}

	if opts.ConnectedPlayer == nil && typeOfBot == human {
		return nil, errors.New("cannot create a bot of human type without a connected Player")
	}

	return &Bot{
		id:        opts.Id,
		name:      opts.Name,
		typeOfBot: typeOfBot,
		player:    opts.ConnectedPlayer,
		helpCount: opts.HelpCount,
	}, nil
}

func (b *Bot) Id() string {
	return b.id
}

func (b *Bot) Name() string {
	return b.name
}

func (b *Bot) IsAi() bool {
	return b.typeOfBot == ai
}

func (b *Bot) IsHuman() bool {
	return b.typeOfBot == human
}

func (b *Bot) CanGetHelp() bool {
	return b.helpCount > 0
}

func (b *Bot) ConnectPlayer(player *Player) error {
	if player == nil {
		return errors.New("Cannot connect an empty player")
	}

	if b.typeOfBot == human && b.player.id != player.id {
		return errors.New("Cannot replace the connected player")
	}

	if b.typeOfBot != ai {
		return errors.New("Can only conect to bot that is currently ai")
	}

	b.typeOfBot = human
	b.player = player
	return nil
}

func RandomBotNames() []string {
	botNames := []string{
		"C-21PO", "R4-D4", "Gart", "HAL 99", "Avis", "ED-I", "Davide", "B.O.B.Z", "T-3PO", "Sort", "T-800X", "EVE-a-L", "GLaDOSE",
	}
	rand.Shuffle(len(botNames), func(i, j int) {
		botNames[i], botNames[j] = botNames[j], botNames[i]
	})
	return botNames[0:totalNumberOfBotsPerGame]
}
