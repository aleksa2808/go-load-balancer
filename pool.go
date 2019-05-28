package main

type pool []*worker

func initPool(nWorker, workerMaxLoad int, done chan *worker) pool {
	pool := make(pool, nWorker)
	for i := range pool {
		w := &worker{
			requests: make(chan Request, workerMaxLoad),
			index:    i,
		}
		pool[i] = w
		go w.work(done)
	}
	return pool
}

func (p pool) Len() int {
	return len(p)
}

func (p pool) Less(i int, j int) bool {
	return p[i].pending < p[j].pending
}

func (p pool) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *pool) Push(x interface{}) {
	*p = append(*p, x.(*worker))
}

func (p *pool) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}
