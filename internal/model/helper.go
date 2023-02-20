package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// This function helps in testing Game Objects
func AssertEqualGame(t *testing.T, expected, actual *Game) {
	assert.Equal(t, expected.id, actual.id, "game id is not equal")
	assert.Equal(t, expected.state, actual.state, "game state is not equal")
	assert.Equal(t, expected.currentTurnIndex, actual.currentTurnIndex, "game currentTurnIndex is not equal")
	assert.Equal(t, expected.turnOrder, actual.turnOrder, "game turnOrder is not equal")
	assert.Equal(t, expected.stateHandled, actual.stateHandled, "game stateHandled is not equal")
	assert.Equal(t, expected.stateTotalTime, actual.stateTotalTime, "game id is not equal")
	assert.EqualValues(t, expected.bots, actual.bots, "game bots is not equal")
	assert.Equal(t, expected.id, actual.id, "game id is not equal")

	// Since we cannot mock postgres time operations. We just check that the updated times are near expected times.
	delta := 5 * time.Second
	assert.WithinRange(t, actual.createdAt, expected.createdAt.Add(-1*delta), expected.createdAt.Add(delta), "game createdAt is not within range")
	assert.WithinRange(t, actual.updatedAt, expected.updatedAt.Add(-1*delta), expected.updatedAt.Add(delta), "game updatedAt is not recent enough")
	if expected.stateHandledAt != nil {
		assert.WithinRange(t, *actual.stateHandledAt, expected.stateHandledAt.Add(-1*delta), expected.stateHandledAt.Add(delta), "game stateHandledAt is not recent enough")
	}
}
