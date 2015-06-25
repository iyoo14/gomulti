package main

import (
	"fmt"
	"github.com/iyoo14/gomulti/entity"
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

type thread_arg struct {
	id    int
	queue *entity.Queue
}

func producer_func(args *thread_arg) {
	for i := 0; i < THREAD_DATA_NUM; i++ {
		num := args.id*THREAD_DATA_NUM + i
		args.queue.Enqueue(num, args.id)
		fmt.Printf("[Producer %d-%d] ==> %d \n", args.id, i, num)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}

	args.queue.Enqueue(END_DATA, args.id)
	return
}

func consumer_func(args *thread_arg, wg *sync.WaitGroup) {
	var i int
	for {
		i = args.queue.Dequeue(args.id)
		if i == END_DATA {
			break
		}
		fmt.Printf("[Consumer %d]     ==> %d \n", args.id, i)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
	wg.Done()
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var ptarg [THREAD_NUM]*thread_arg
	var ctarg [THREAD_NUM]*thread_arg

	var q *entity.Queue = entity.NewQueue(MAX_QUEUE_NUM)
	var wg = new(sync.WaitGroup)

	for i := 0; i < THREAD_NUM; i++ {
		ptarg[i] = &thread_arg{i, q}
		go producer_func(ptarg[i])
	}

	for i := 0; i < THREAD_NUM; i++ {
		ctarg[i] = &thread_arg{i, q}
		wg.Add(1)
		go consumer_func(ctarg[i], wg)
	}
	wg.Wait()
	//time.Sleep(time.Second * 100)
}
