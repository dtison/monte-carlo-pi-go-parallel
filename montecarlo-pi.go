package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const MaxUint = ^uint32(0)
const MinUint = 0
const MaxInt = int32(MaxUint >> 1)
const MinInt = -MaxInt - 1

func main() {

	start := time.Now()

	fmt.Println("GOMAXPROCS set from", runtime.GOMAXPROCS(runtime.NumCPU()), "to", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().UnixNano())
	monteCarloPi(1000000000)

	elapsed := time.Since(start)
	fmt.Printf("Processing finished in %f seconds\n", float32(elapsed)/float32(time.Millisecond)/1000)
}

func monteCarloPi(samples int) {

	numCPUs := runtime.NumCPU()
	samplesPerThread := samples / numCPUs
	threadResults := make(chan uint64, numCPUs)

	ticker := time.NewTicker(500 * time.Millisecond)

	for i := 0; i < numCPUs; i++ {

		go func(cpu int) {

			var pointsInside uint64
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for j := 0; j < samplesPerThread; j++ {

				x, y := r.Float64(), r.Float64()

				if x*x+y*y <= 1.0 {
					pointsInside++
				}

				if j&0x3ffff == 262143 && cpu == 0 {
					select {
					// case t := <-ticker.C:
					// 	fmt.Println("Tick at", t)
					case <-ticker.C:
						fmt.Printf("%f\n", float32(j)/float32(samplesPerThread))
					default:
					}
				}

			}
			threadResults <- pointsInside
		}(i)

	}

	var total uint64
	for i := 0; i < numCPUs; i++ {
		total += <-threadResults
	}

	fmt.Printf("%f\n", float64(total)/float64(samples)*4.0)

}

//	updateInterval := samplesPerThread >> 7

// if cpu == 0 && (j&0xffff) == MaxInt  {
// 	fmt.Printf(".")

// }

// if cpu == 0 && (j%updateInterval) == 0 {
// 	fmt.Printf(".")
// }
