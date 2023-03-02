package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type Player struct {
	id string
}

type PlayerOptions struct {
	Id string
}

func NewPlayer(opts PlayerOptions) (*Player, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.Errorf("cannot create player with an empty id")
	}
	return &Player{
		id: opts.Id,
	}, nil
}

func (p *Player) Id() string {
	return p.id
}
