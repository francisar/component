package pool


type Pool struct {
	work chan interface{}   // 任务
	sem  chan struct{} // 数量
	handler func(interface{})
}


func NewPool(size int, handler func(interface{})) *Pool {
	return &Pool{
		work: make(chan interface{}),
		sem:  make(chan struct{}, size),
		handler: handler,
	}
}

func (p *Pool) NewTask(task interface{}) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *Pool) worker(task interface{}) {
	defer func() { <-p.sem }()
	for {
		p.handler(task)
		task = <-p.work
	}
}