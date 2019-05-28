package main

type worker struct {
	requests chan Request
	pending  int
	index    int
}

func (w *worker) work(done chan *worker) {
	for req := range w.requests {
		req.r <- req.fn()
		done <- w
	}
}
