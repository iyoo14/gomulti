package main

import (
	"fmt"
	"github.com/iyoo14/gomulti/entity"
	"math"
	"runtime"
)

const (
	THREAD_NUM      = 2
	DATA_NUM        = 100
	THREAD_DATA_NUM = DATA_NUM / THREAD_NUM
	END_DATA        = -1
	QUEUE_NUM       = DATA_NUM / 2
)

type thread_arg struct {
	divider int
	prev_q  *entity.Queue
	next_q  *entity.Queue
}

func master_func(args *thread_arg) {
	for i := 2; i < DATA_NUM; i++ {
		if i%2 != 0 || i == 2 {
			args.next_q.Enqueue(i)
		}
	}
	args.next_q.Enqueue(END_DATA)
}

func thread_func(args *thread_arg) {
	for i := 2; i < DATA_NUM+1; i++ {
		v := args.prev_q.Dequeue()
		if v == END_DATA {
			args.next_q.Enqueue(v)
			break
		}
		if v%args.divider != 0 || args.divider > int(math.Sqrt(float64(v))) {
			args.next_q.Enqueue(v)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var pipeline_size int = (int(math.Sqrt(float64(DATA_NUM))) + 1) / 2
	var targ []*thread_arg = make([]*thread_arg, pipeline_size)
	var qs []*entity.Queue = make([]*entity.Queue, pipeline_size)

	for i := 0; i < pipeline_size; i++ {
		qs[i] = entity.NewQueue(DATA_NUM + 1)
	}

	targ[0] = &thread_arg{divider: 2, next_q: qs[0]}
	go master_func(targ[0])
	for i := 1; i < pipeline_size; i++ {
		targ[i] = &thread_arg{divider: i*2 + 1, prev_q: qs[i-1], next_q: qs[i]}
		go thread_func(targ[i])
	}

	fmt.Println(pipeline_size - 1)
	for i := 2; i < DATA_NUM+1; i++ {
		v := qs[pipeline_size-1].Dequeue()

		if v == END_DATA {
			break
		}

		fmt.Printf("%d ", v)
	}
	fmt.Printf("\n")
}
