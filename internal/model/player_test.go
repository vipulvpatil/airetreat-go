package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPlayer(t *testing.T) {
	tests := []struct {
		name           string
		input          PlayerOptions
		expectedOutput *Player
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          PlayerOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create player with an empty id",
		},
		{
			name: "Player gets created successfully",
			input: PlayerOptions{
				Id: "123",
			},
			expectedOutput: &Player{
				id: "123",
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPlayer(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Player_Id(t *testing.T) {
	tests := []struct {
		name           string
		input          *Player
		expectedOutput string
	}{
		{
			name:           "returns empty if player is empty",
			input:          nil,
			expectedOutput: "",
		},
		{
			name:           "returns Id successfully",
			input:          &Player{id: "id1"},
			expectedOutput: "id1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Id()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
