package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type Bot struct {
	id string
}

type BotOptions struct {
	Id string
}

func NewBot(opts BotOptions) (*Bot, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create bot with an empty id")
	}
	return &Bot{id: opts.Id}, nil
}
