package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

const MaxUint = ^uint32(0)
const MinUint = 0
const MaxInt = int32(MaxUint >> 1)
const MinInt = -MaxInt - 1

var startTime time.Time

func main() {

	startTime = time.Now()

	fmt.Println("Using", runtime.NumCPU(), "cores..\n")

	samples := flag.Uint64("samples", 1000000000, "Number of samples to test")
	flag.Parse()

	monteCarloPi(*samples)

	displayMessageWithElapsedTime("Processing finished.")
}

func monteCarloPi(samples uint64) {

	numCPUs := runtime.NumCPU()
	samplesPerThread := uint64(samples / uint64(numCPUs))

	samplesCoreZero := uint64(float64(samplesPerThread) * .7)
	samplesPerThread = uint64(math.Round(float64(samples-samplesCoreZero) / float64((numCPUs - 1))))

	totalSamples := samplesPerThread*(uint64(numCPUs-1)) + samplesCoreZero

	fmt.Printf("Corezero samples %d, Rest of threads %d total samples %d\n", samplesCoreZero, samplesPerThread, totalSamples)

	threadResults := make(chan uint64, numCPUs)

	ticker := time.NewTicker(100 * time.Millisecond)

	bar := pb.StartNew(100)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	go threadMCUI(samplesCoreZero, threadResults, ticker, bar, &wg)

	for i := 0; i < numCPUs; i++ {
		go threadMC(samplesPerThread, threadResults, &wg)
	}
	wg.Wait()

	var total uint64

	for i := 0; i < numCPUs; i++ {
		total += <-threadResults
	}

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
	defer wg.Done()
	defer bar.Finish()

	//	defer displayMessageWithElapsedTime("UI Thread now finished.")
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
}

func displayMessageWithElapsedTime(message string) {
	elapsed := time.Since(startTime)
	fmt.Printf("%s %f seconds\n", message, float32(elapsed)/float32(time.Millisecond)/1000)
}

// for result := range threadResults {
// 	total += result
// }
