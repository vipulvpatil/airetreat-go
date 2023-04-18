package storage

import "github.com/pkg/errors"

type BotAccessorMockSuccess struct {
}

func (p *BotAccessorMockSuccess) UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error {
	return nil
}

func (p *BotAccessorMockSuccess) UpdateBotDecrementHelpCountUsingTransaction(botId string, transaction DatabaseTransaction) error {
	return nil
}

type BotAccessorMockFailure struct {
}

func (p *BotAccessorMockFailure) UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error {
	return errors.New("unable to update bot")
}

func (p *BotAccessorMockFailure) UpdateBotDecrementHelpCountUsingTransaction(botId string, transaction DatabaseTransaction) error {
	return errors.New("unable to update bot")
}
