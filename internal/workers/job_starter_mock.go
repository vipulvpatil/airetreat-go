package workers

import (
	"github.com/pkg/errors"

	"github.com/gocraft/work"
)

type JobStarterMockSuccess struct{}

func (j *JobStarterMockSuccess) Enqueue(jobName string, args map[string]interface{}) (*work.Job, error) {
	return &work.Job{}, nil
}

type JobStarterMockFailure struct{}

func (j *JobStarterMockFailure) Enqueue(jobName string, args map[string]interface{}) (*work.Job, error) {
	return nil, errors.New("unable to enqueue job")
}
