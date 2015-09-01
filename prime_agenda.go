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

func judgmentPrime() {

	for {
		q.Dequeue(
