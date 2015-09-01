package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

const (
	THREAD_NUM = 2
	DATA_NUM   = 10
)

type thread_arg struct {
	id     int
	primes []bool
	wg     *sync.WaitGroup
}

func thread_func(args *thread_arg) {

	rng := (DATA_NUM-2)/THREAD_NUM + 1
	c_start := 2 + args.id*rng
	c_end := 2 + (args.id+1)*rng
	if c_end > DATA_NUM {
		c_end = DATA_NUM
	}
	for i := c_start; i < c_end; i++ {
		limit := int(math.Sqrt(float64(i)))
		for j := 2; j <= limit; j++ {
			if args.primes[j] && i%j == 0 {
				args.primes[i] = false
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

	for i := 0; i < THREAD_NUM; i++ {
		targ[i] = &thread_arg{i, primes[:], wg}
		wg.Add(1)
		go thread_func(targ[i])
	}
	wg.Wait()
	for k, v := range primes {
		fmt.Printf("%d - %v\n", k, v)
	}
}
