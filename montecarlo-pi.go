package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

/* type count64 int64

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
} */

func main() {
	fmt.Println("GOMAXPROCS set from", runtime.GOMAXPROCS(runtime.NumCPU()), "to", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().UnixNano())
	monteCarloPi(5000000000)
}

func monteCarloPi(samples int) {

	numCPUs := runtime.NumCPU()

	samplesPerThread := samples / numCPUs

	threadResults := make(chan uint64, numCPUs)

	for i := 0; i < numCPUs; i++ {
		go func() {
			var pointsInside uint64
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for j := 0; j < samplesPerThread; j++ {

				x, y := r.Float64(), r.Float64()

				if x*x+y*y <= 1.0 {
					pointsInside++
				}
			}
			threadResults <- pointsInside
		}()

		// x, y := rand.Float64(), rand.Float64()
		// if (x*x)+(y*y) <= 1.0 {
		// 	pointsInside++
		// }

	}
	//	fmt.Println(4.0 * float64(<-points2) / float64(n))

	var total uint64
	// for result := range threadResults {
	// 	total += result
	// }

	for i := 0; i < numCPUs; i++ {
		total += <-threadResults
	}
	//total *= 4.0

	fmt.Printf("%f\n", float64(total)/float64(samples)*4.0)

}
