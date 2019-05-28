package main

import (
	"container/heap"
)

type balancer struct {
	pool  pool
	done  chan *worker
	count int
}

// Balance distributes work over a specified number of workers.
func Balance(work chan Request, nWorker, workerMaxLoad int) {
	done := make(chan *worker, nWorker*workerMaxLoad)
	pool := initPool(nWorker, workerMaxLoad, done)

	b := &balancer{
		pool: pool,
		done: done,
	}

	for {
		// // simple visualisation for
		// // a small number of workers
		// for i := 0; i < 20; i++ {
		// 	fmt.Println()
		// }
		// for id, w := range pool {
		// 	fmt.Printf("%d: ", id)
		// 	for i := 0; i < workerMaxLoad; i++ {
		// 		if i < w.pending {
		// 			fmt.Printf("%d ", i)
		// 		} else {
		// 			fmt.Print("_ ")
		// 		}
		// 	}
		// 	fmt.Print("\n")
		// }

		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *balancer) dispatch(r Request) {
	w := b.pool[0]
	select {
	case w.requests <- r:
		w.pending++
		heap.Fix(&b.pool, 0)
	default:
		// servers full, discard request
		r.r <- -1
	}
}

func (b *balancer) completed(w *worker) {
	w.pending--
	heap.Fix(&b.pool, w.index)
	b.count++
}
