package main

// Request specifies a work function
// and a result channel.
type Request struct {
	fn func() int
	r  chan int
}
