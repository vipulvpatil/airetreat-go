package storage

import "errors"

type MessageCreatorMockSuccess struct {
	PlayerId string
}

func (m *MessageCreatorMockSuccess) CreateMessage(sourceBotId, targetBotId, text, messageType string) error {
	return nil
}

func (m *MessageCreatorMockSuccess) CreateMessageUsingTransaction(sourceBotId, targetBotId, text, messageType string, transation DatabaseTransaction) error {
	return nil
}

type MessageCreatorMockFailure struct {
}

func (m *MessageCreatorMockFailure) CreateMessage(sourceBotId, targetBotId, text, messageType string) error {
	return errors.New("unable to create message")
}

func (m *MessageCreatorMockFailure) CreateMessageUsingTransaction(sourceBotId, targetBotId, text, messageType string, transation DatabaseTransaction) error {
	return errors.New("unable to create message")
}
