# monte-carlo-pi-go-parallel
**Pi Estimator using Monte Carlo Algorithm with Parallel Goroutines**

**Options**

  -samples=number of samples

**Packages**

  go get github.com/schollz/progressbar

**Overview**

For speed and accuracy counters and loops all use uint64.

Something to try would be using the big int package so as to enable huge numbers of samples.

Computational work is divided across the number of available compute cores, and processing runs in parallel.

One core is assigned UI work, and is allocated fewer samples than the other cores.

A background timer is used to manage UI progress bar updating.







