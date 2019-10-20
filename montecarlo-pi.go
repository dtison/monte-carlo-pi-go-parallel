package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

/*
const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1


} */

func main() {
	fmt.Println("GOMAXPROCS set from", runtime.GOMAXPROCS(runtime.NumCPU()), "to", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().UnixNano())
	monteCarloPi(1000000000)
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

	}

	var total uint64
	for i := 0; i < numCPUs; i++ {
		total += <-threadResults
	}

	fmt.Printf("%f\n", float64(total)/float64(samples)*4.0)

}
