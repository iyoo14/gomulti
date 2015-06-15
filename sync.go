package main

import (
	"fmt"
	"math/rand"
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
	values    [MAX_QUEUE_NUM]int
	remain    int
	rp        int
	wp        int
	mutex     *sync.Mutex
	not_full  *sync.Cond
	not_empty *sync.Cond
}

type thread_arg struct {
	id    int
	queue *queue
}

func enqueue(q *queue, v int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for q.remain == MAX_QUEUE_NUM {
		q.not_full.Wait()
	}
	q.values[q.wp] = v
	q.wp++
	q.remain++
	if q.wp == MAX_QUEUE_NUM {
		q.wp = 0
	}
	q.not_empty.Signal()
}

func dequeue(q *queue, v *int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for q.remain == 0 {
		q.not_empty.Wait()
	}
	*v = q.values[q.rp]
	q.rp++
	q.remain--
	if q.rp == MAX_QUEUE_NUM {
		q.rp = 0
	}
	q.not_full.Signal()
}

func producer_func(args *thread_arg) {
	for i := 0; i < THREAD_DATA_NUM; i++ {
		num := args.id*THREAD_DATA_NUM + i
		enqueue(args.queue, num)
		fmt.Printf("[Producer %d] ==> %d \n", args.id, num)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}

	enqueue(args.queue, END_DATA)
	return
}

func consumer_func(args *thread_arg) {
	var i int
	for {
		dequeue(args.queue, &i)
		if i == END_DATA {
			break
		}
		fmt.Printf("[Consumer %d]     ==> %d \n", i)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
	return
}

func main() {
	var ptarg [THREAD_NUM]*thread_arg
	var ctarg [THREAD_NUM]*thread_arg

	var mutex = new(sync.Mutex)
	var nf = new(sync.Mutex)
	var not_full = sync.NewCond(nf)
	var ne = new(sync.Mutex)
	var not_empty = sync.NewCond(ne)
	var q queue = queue{}

	q.rp = 0
	q.wp = 0
	q.remain = 0
	q.mutex = mutex
	q.not_full = not_full
	q.not_empty = not_empty

	for i := 0; i < THREAD_NUM; i++ {
		ptarg[i] = &thread_arg{i, &q}
		go producer_func(ptarg[i])
	}

	for i := 0; i < THREAD_NUM; i++ {
		ctarg[i] = &thread_arg{i, &q}
		go consumer_func(ctarg[i])
	}

	for i := 0; i < THREAD_NUM; i++ {

	}
	time.Sleep(time.Second * 100)
}
