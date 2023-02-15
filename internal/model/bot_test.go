package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBot(t *testing.T) {
	tests := []struct {
		name           string
		input          BotOptions
		expectedOutput *Bot
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          BotOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create bot with an empty id",
		},
		{
			name:           "name is empty",
			input:          BotOptions{Id: "1"},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create bot with an empty name",
		},
		{
			name: "Bot gets created successfully",
			input: BotOptions{
				Id:   "123",
				Name: "some name",
			},
			expectedOutput: &Bot{
				id:   "123",
				name: "some name",
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBot(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_RandomBotNames(t *testing.T) {
	tests := []struct {
		name           string
		input          int64
		expectedOutput []string
	}{
		{
			name:           "random names are generated",
			input:          10,
			expectedOutput: []string{"The Hivey-five", "Gart", "RoboCupp", "T-3PO", "Avis"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomBotNames(tt.input)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
