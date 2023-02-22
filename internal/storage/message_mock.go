package storage

import "errors"

type MessageCreatorMockSuccess struct {
	PlayerId string
}

func (m *MessageCreatorMockSuccess) CreateMessage(botId, text string) error {
	return nil
}

type MessageCreatorMockFailure struct {
}

func (m *MessageCreatorMockFailure) CreateMessage(botId, text string) error {
	return errors.New("unable to create message")
}
