package worker_pool

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string
const (
	Pending JobStatus = "pending"
	InQueue JobStatus = "in-queue"
	Running JobStatus = "running"
	Suceed JobStatus = "suceed"
	Failed JobStatus = "failed"
)

type Job struct {
	ID     uuid.UUID
	Status JobStatus
}

func NewJob()*Job{
	return &Job{
		ID: uuid.New(),
		Status: Pending,
	}
}

func (j *Job) Perform() {
	j.Status = Running

	time.Sleep(3 * time.Second)

	j.Status = Suceed
}
