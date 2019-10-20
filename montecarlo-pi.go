package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

type count64 int64

func (c *count64) increment() int64 {
	var next int64
	for {
		next = int64(*c) + 1
		if atomic.CompareAndSwapInt64((*int64)(c), int64(*c), next) {
			return next
		}
	}
}
func (c *count64) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func main() {
	fmt.Println("GOMAXPROCS set from", runtime.GOMAXPROCS(runtime.NumCPU()), "to", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().UnixNano())
	monteCarloPi(1000)
}

func monteCarloPi(n int) {
	var pointsInside uint64
	points2 := make(chan uint64)

	//	var wg sync.WaitGroup
	//	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			x, y := rand.Float64(), rand.Float64()

			if (x*x)+(y*y) <= 1.0 {
				//pointsInside.increment()
				// func AddUint64(addr *uint64, delta uint64) (new uint64)
				//atomic.AddUint64(&pointsInside, 1)

				points2 <- 1
			}
			//			wg.Done()
		}()

		select {
		case <-points2:
			pointsInside++
		case i == n:

		}

		// x, y := rand.Float64(), rand.Float64()
		// if (x*x)+(y*y) <= 1.0 {
		// 	pointsInside++
		// }

	}
	fmt.Println(4.0 * float64(<-points2) / float64(n))

	//	wg.Wait()
}
