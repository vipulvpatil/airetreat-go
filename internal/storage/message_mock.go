package storage

import "errors"

type MessageCreatorMockSuccess struct {
	PlayerId string
}

func (m *MessageCreatorMockSuccess) AddMessage(botId, text string) error {
	return nil
}

type MessageCreatorMockFailure struct {
}

func (m *MessageCreatorMockFailure) AddMessage(botId, text string) error {
	return errors.New("unable to create message")
}
