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
	values    [MAX_QUEUE_NUM]int
	remain    int
	rp        int
	wp        int
	wg        *sync.WaitGroup
	not_full  *sync.Cond
	not_empty *sync.Cond
}

type thread_arg struct {
	id    int
	queue *queue
}

func enqueue(q *queue, v int, id int) {
	fmt.Printf("en wait %d\n", id)
	q.not_full.L.Lock()
	fmt.Printf("en start %d\n", id)
	defer q.not_full.L.Unlock()
	for q.remain == MAX_QUEUE_NUM {
		q.not_full.Wait()
	}
	q.values[q.wp] = v
	q.wp++
	q.remain++
	if q.wp == MAX_QUEUE_NUM {
		q.wp = 0
	}
	//q.not_empty.Signal()
	q.not_full.Signal()
	fmt.Printf("en end %d\n", id)
}

func dequeue(q *queue, v *int, id int) {
	//q.not_empty.L.Lock()
	fmt.Printf("de wait %d\n", id)
	q.not_full.L.Lock()
	fmt.Printf("de start %d\n", id)
	//defer q.not_empty.L.Unlock()
	defer q.not_full.L.Unlock()
	for q.remain == 0 {
		fmt.Printf("q empty %d\n", id)
		//q.not_empty.Wait()
		q.not_full.Wait()
	}
	*v = q.values[q.rp]
	q.rp++
	q.remain--
	if q.rp == MAX_QUEUE_NUM {
		q.rp = 0
	}
	q.not_full.Signal()
	fmt.Printf("de end %d\n", id)
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
		dequeue(args.queue, &i, args.id)
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
	var not_full = sync.NewCond(lock)
	var not_empty = sync.NewCond(lock)
	var q queue = queue{}
	var wg = new(sync.WaitGroup)

	q.rp = 0
	q.wp = 0
	q.remain = 0
	q.wg = wg
	q.not_full = not_full
	q.not_empty = not_empty

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
