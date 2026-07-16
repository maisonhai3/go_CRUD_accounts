package worker_pool

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
	p.q <- j
	j.Status = InQueue
}

func (p *Producer) CloseJobQueue(){
	close(p.q)
}