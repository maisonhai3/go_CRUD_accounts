package phase1semaphores

type Producer struct {
	q JobQueue
	jobs []*Job
}

func (p *Producer) PushJob(){
}

func (p *Producer) CloseJobQueue(){
}