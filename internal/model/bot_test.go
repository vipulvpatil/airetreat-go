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
			name: "Bot gets created successfully",
			input: BotOptions{
				Id: "123",
			},
			expectedOutput: &Bot{
				id: "123",
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
