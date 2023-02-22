package storage

import "errors"

type PlayerCreatorMockSuccess struct {
	PlayerId string
}

func (p *PlayerCreatorMockSuccess) CreatePlayer() (string, error) {
	return p.PlayerId, nil
}

type PlayerCreatorMockFailure struct {
}

func (p *PlayerCreatorMockFailure) CreatePlayer() (string, error) {
	return "", errors.New("unable to create player")
}
