package worker_pool

import "fmt"

type Producer struct {
	q chan *Job
	jobs []*Job
}

func NewProducer(q chan *Job) *Producer{
	return &Producer{
		q: q,
	}
}

func (p *Producer) PushJob(j *Job){
	j.Status = InQueue
	fmt.Printf("Job %v is pushed to queue", j.ID)
	p.q <- j
}

func (p *Producer) CloseJobQueue(){
	close(p.q)
}