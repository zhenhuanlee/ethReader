package ppool

import "fmt"

// func main() {
// 	pool := NewPool(1)

// 	a := "1"

// 	for {
// 		time.Sleep(1000 * time.Millisecond)
// 	}
// }

// Pool struct
type Pool struct {
	Count   int
	JobChan chan func()
	workers *[]*worker
}

type worker struct {
	name  string
	state string
}

// NewPool init pool
func NewPool(count int) *Pool {
	pool := &Pool{
		Count:   count,
		JobChan: make(chan func()),
	}
	pool.workers = pool.initWorkers()
	return pool
}

func (pool *Pool) initWorkers() *[]*worker {
	workers := make([]*worker, pool.Count)

	for i := 0; i < pool.Count; i++ {
		w := &worker{
			name:  fmt.Sprintf("worker-%d", i),
			state: "waiting",
		}

		workers[i] = w.run(pool)
	}

	return &workers
}

func (w worker) run(pool *Pool) *worker {
	go func() {
		for j := range pool.JobChan {
			fmt.Printf("%s is working\n", w.name)
			w.state = "working"
			j()
			w.state = "waitting"
			fmt.Printf("%s done\n", w.name)
		}
	}()
	return &w
}
