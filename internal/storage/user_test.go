package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
)

func Test_UserByEmail(t *testing.T) {
	returnUser, _ := model.NewUser(model.UserOptions{
		Id:    "123",
		Email: "test@example.com",
	})
	tests := []struct {
		name            string
		input           string
		expectedOutput  *model.User
		setupSqlStmts   []string
		cleanupSqlStmts []string
		errorExpected   bool
		errorString     string
	}{
		{
			name:           "returns user successfully",
			input:          "test@example.com",
			expectedOutput: returnUser,
			setupSqlStmts: []string{
				`INSERT INTO public."User" (id, email) VALUES ('123', 'test@example.com')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."User" WHERE email = 'test@example.com'`,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:            "errors if user not in db",
			input:           "test@example.com",
			expectedOutput:  returnUser,
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			errorExpected:   true,
			errorString:     "UserByEmail test@example.com: no such user",
		},
		{
			name:            "errors if blank email provided",
			input:           "   ",
			expectedOutput:  returnUser,
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			errorExpected:   true,
			errorString:     "cannot search by blank email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db: testDb,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)

			user, err := s.UserByEmail(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, user, tt.expectedOutput)
			} else {
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
