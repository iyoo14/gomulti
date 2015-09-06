package entity

import (
	"sync"
	"time"
)

type Queue struct {
	values chan int
	lock   *sync.Mutex
}

func (q *Queue) Enqueue(v int) {
	q.values <- v
}

func (q *Queue) Dequeue() (v int) {
	/*
		q.lock.Lock()
		defer q.lock.Unlock()
		for {
			if len(q.values) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	*/
	time.Sleep(1 * time.Second)
	return <-q.values
}

func NewQueue(num int) *Queue {
	q := &Queue{}
	q.values = make(chan int, num)
	q.lock = new(sync.Mutex)
	return q
}
