package main

import (
	"fmt"
	"math/rand"
	"time"
)

func requester(id int, work chan<- Request) {
	result := make(chan int)
	for {
		// wait a bit
		time.Sleep(100 * time.Millisecond)

		// send a request
		work <- Request{
			fn: func() int {
				t := 500 + rand.Intn(501)
				time.Sleep(time.Duration(t) * time.Millisecond)
				return t
			},
			r: result,
		}

		// wait for result
		res := <-result

		// announce if request got denied
		if res == -1 {
			fmt.Printf("R%d:\trequest unhandled\n", id)
		}
	}
}

const (
	nRequestors   = 20
	nWorkers      = 5
	workerMaxLoad = 10
)

func main() {
	// create main channel on which work is passed
	work := make(chan Request, 2*nRequestors)

	// start dummy requestors
	for i := 0; i < nRequestors; i++ {
		go requester(i, work)
	}

	// start balancing the work channel
	Balance(work, nWorkers, workerMaxLoad)
}
