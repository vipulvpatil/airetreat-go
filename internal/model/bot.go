package model

import (
	"math/rand"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

const totalNumberOfBotsPerGame = 5

type Bot struct {
	id   string
	name string
}

type BotOptions struct {
	Id   string
	Name string
}

func NewBot(opts BotOptions) (*Bot, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create bot with an empty id")
	}
	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create bot with an empty name")
	}
	return &Bot{
		id:   opts.Id,
		name: opts.Name,
	}, nil
}

func RandomBotNames(randomSeed int64) []string {
	botNames := []string{
		"C-21PO", "R4-D4", "Gart", "HAL 9999", "Avis", "ED-I", "T-5000", "Davide", "B.O.B.Z", "The Machy-ne", "GLaDOODLES", "JARV-EESE", "The Hivey-five", "T-3PO", "InfoData", "Sort", "Electronic Device-209", "T-800X", "RoboCupp", "EVE-a-L", "GLaDOSE",
	}
	rand.Seed(randomSeed)
	rand.Shuffle(len(botNames), func(i, j int) {
		botNames[i], botNames[j] = botNames[j], botNames[i]
	})
	return botNames[0:totalNumberOfBotsPerGame]
}
