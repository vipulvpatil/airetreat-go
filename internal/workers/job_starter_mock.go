package workers

import (
	"github.com/pkg/errors"

	"github.com/gocraft/work"
)

type JobStarterMockCallCheck struct {
	CalledArgs []map[string]interface{}
}

func (j *JobStarterMockCallCheck) EnqueueUnique(jobName string, args map[string]interface{}) (*work.Job, error) {
	if j.CalledArgs == nil {
		j.CalledArgs = []map[string]interface{}{}
	}
	j.CalledArgs = append(j.CalledArgs, args)
	return &work.Job{}, nil
}

type JobStarterMockFailure struct{}

func (j *JobStarterMockFailure) Enqueue(jobName string, args map[string]interface{}) (*work.Job, error) {
	return nil, errors.New("unable to enqueue job")
}
