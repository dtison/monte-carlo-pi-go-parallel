package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

// TODO:  samplesPerThread should be uint64

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

func monteCarloPi(samples uint64) {

	numCPUs := runtime.NumCPU()
	samplesPerThread := uint64(samples / uint64(numCPUs))

	samplesCoreZero := samplesPerThread >> 1
	samplesPerThread = uint64(math.Round(float64(samples-samplesCoreZero) / float64((numCPUs - 1))))

	totalSamples := samplesPerThread*(uint64(numCPUs-1)) + samplesCoreZero

	fmt.Printf("Corezero samples %d, Rest of threads %d total samples %d\n", samplesCoreZero, samplesPerThread, totalSamples)

	threadResults := make(chan uint64, numCPUs)

	ticker := time.NewTicker(10 * time.Millisecond)

	bar := pb.StartNew(100)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	go threadMCUI(samplesCoreZero, threadResults, ticker, bar, &wg)

	for i := 0; i < numCPUs; i++ {
		go threadMC(samplesPerThread, threadResults, &wg)
	}
	wg.Wait()
	// TODO:  Try the range for again

	var total uint64
	for i := 0; i < numCPUs; i++ {
		total += <-threadResults
	}
	bar.Finish()

	fmt.Printf("%f\n", float64(total)/float64(totalSamples)*4.0)

}

func threadMC(samples uint64, threadResults chan uint64, wg *sync.WaitGroup) {
	defer wg.Done()

	var pointsInside uint64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for j := uint64(0); j < samples; j++ {

		x, y := r.Float64(), r.Float64()

		if x*x+y*y <= 1.0 {
			pointsInside++
		}

	}

	threadResults <- pointsInside
}

func threadMCUI(samples uint64, threadResults chan uint64, ticker *time.Ticker, bar *pb.ProgressBar, wg *sync.WaitGroup) {

	var pointsInside uint64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	uiPos := uint64(0)
	uiAdvance := samples / 100
	for j := uint64(0); j < samples; j++ {

		x, y := r.Float64(), r.Float64()

		if x*x+y*y <= 1.0 {
			pointsInside++
		}

		select {
		case <-ticker.C:
			for i := (uiPos * uiAdvance); i < j; i += uiAdvance {
				bar.Increment()
				uiPos++

			}
		default:
		}

	}
	threadResults <- pointsInside
	defer wg.Done()
}

/*
func testPoints(r *rand.Rand, pointsInside *uint64) {
	x, y := r.Float64(), r.Float64()

	if x*x+y*y <= 1.0 {
		*pointsInside++
	}
} */

//	updateInterval := samplesPerThread >> 7

// if cpu == 0 && (j&0xffff) == MaxInt  {
// 	fmt.Printf(".")

// }

// if cpu == 0 && (j%updateInterval) == 0 {
// 	fmt.Printf(".")
// }

/* 				if cpu == 0 {
	   fmt.Printf("j&0xff is %d\n", j&0xff)
   }
   if j&0xff == 255 && cpu == 0 {
	   select {
	   // case t := <-ticker.C:
	   // 	fmt.Println("Tick at", t)
	   case <-ticker.C:
		   //	fmt.Printf("%f\n", float32(j)/float32(samplesPerThread))
		   bar.Increment()
	   default:
	   }
   } */
