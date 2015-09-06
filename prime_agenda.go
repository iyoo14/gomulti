package main

import (
	"fmt"
	"github.com/iyoo14/gomulti/entity"
	"math"
	"runtime"
	"sync"
)

const (
	THREAD_NUM      = 2
	DATA_NUM        = 10
	THREAD_DATA_NUM = DATA_NUM / THREAD_NUM
	END_DATA        = -1
	MAX_QUEUE_SIZE  = 5
)

type thread_arg struct {
	id     int
	primes []bool
	queue  *entity.Queue
	wg     *sync.WaitGroup
}

func thread_func(args *thread_arg) {

	for {
		num := args.queue.Dequeue()
		if num == END_DATA {
			break
		}
		limit := int(math.Sqrt(float64(num)))
		fmt.Printf("id num : %d %d\n", args.id, num)
		for i := 2; i <= limit; i++ {
			if args.primes[i] && num%i == 0 {
				args.primes[num] = false
				break
			}
		}
	}
	args.wg.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var targ [THREAD_NUM]*thread_arg
	var wg = new(sync.WaitGroup)
	var primes [DATA_NUM]bool
	for i := 0; i < DATA_NUM; i++ {
		primes[i] = true
	}

	var q *entity.Queue = entity.NewQueue(MAX_QUEUE_SIZE)

	for i := 0; i < THREAD_NUM; i++ {
		targ[i] = &thread_arg{i, primes[:], q, wg}
		wg.Add(1)
		go thread_func(targ[i])
	}

	for i := 2; i < DATA_NUM; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < THREAD_NUM; i++ {
		q.Enqueue(END_DATA)
	}

	wg.Wait()
	for k, v := range primes {
		fmt.Printf("%d - %v\n", k, v)
	}
}
