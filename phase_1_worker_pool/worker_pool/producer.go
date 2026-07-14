package worker_pool

type Producer struct {
	q JobQueue
	jobs []*Job
}

func (p *Producer) PushJob(){
}

func (p *Producer) CloseJobQueue(){
}