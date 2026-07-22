package worker_pool

import (
	"time"
	"github.com/google/uuid"
	"fmt"
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
	uuid := uuid.New()
	fmt.Printf("Job %v is generated", uuid)
	return &Job{
		ID: uuid,
		Status: Pending,
	}
}

func (j *Job) Perform() {
	j.Status = Running

	fmt.Printf("%v started", j.ID.String())
	time.Sleep(3 * time.Second)
	fmt.Printf("%v ended", j.ID.String())

	j.Status = Suceed
}
