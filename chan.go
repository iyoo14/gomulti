package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	THREAD_NUM      = 2
	DATA_NUM        = 10
	MAX_QUEUE_NUM   = 3
	THREAD_DATA_NUM = DATA_NUM / THREAD_NUM
	END_DATA        = -1
)

type queue struct {
	values chan int
	wg     *sync.WaitGroup
	lock   *sync.Mutex
}

type thread_arg struct {
	id    int
	queue *queue
}

func enqueue(q *queue, v int, id int) {
	q.values <- v
}

func dequeue(q *queue, id int) (v int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		if len(q.values) != 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return <-q.values
}

func producer_func(args *thread_arg) {
	for i := 0; i < THREAD_DATA_NUM; i++ {
		num := args.id*THREAD_DATA_NUM + i
		enqueue(args.queue, num, args.id)
		fmt.Printf("[Producer %d] ==> %d \n", args.id, num)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}

	enqueue(args.queue, END_DATA, args.id)
	return
}

func consumer_func(args *thread_arg) {
	var i int
	for {
		i = dequeue(args.queue, args.id)
		if i == END_DATA {
			break
		}
		fmt.Printf("[Consumer %d]     ==> %d \n", args.id, i)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
	args.queue.wg.Done()
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var ptarg [THREAD_NUM]*thread_arg
	var ctarg [THREAD_NUM]*thread_arg

	var lock = new(sync.Mutex)
	var q queue = queue{}
	var wg = new(sync.WaitGroup)

	q.lock = lock
	q.values = make(chan int, MAX_QUEUE_NUM)
	q.wg = wg

	for i := 0; i < THREAD_NUM; i++ {
		ptarg[i] = &thread_arg{i, &q}
		go producer_func(ptarg[i])
	}

	for i := 0; i < THREAD_NUM; i++ {
		ctarg[i] = &thread_arg{i, &q}
		wg.Add(1)
		go consumer_func(ctarg[i])
	}
	wg.Wait()
	//time.Sleep(time.Second * 100)
}
