package main

import (
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"
)

func main() {
	numWorkersFlag := flag.Int("workers", DEFAULT_NUMBER_WORKERS, "Number of worker goroutines to use")
	maxNumberFlag := flag.Int("max", DEFAULT_MAX_NUMBER, "Upper limit for finding prime numbers")
	flag.Parse()

	numWorkers := *numWorkersFlag
	maxNum := *maxNumberFlag

	startTime := time.Now()
	fmt.Printf("Starting to find prime numbers between 1 and %d using %d workers\n", maxNum, numWorkers)

	// The results channel needs larger buffer to prevent blocking when many primes are found
	results := make(chan int, numWorkers*100)
	progress := make(chan ProgressUpdate, numWorkers*100)

	var wg sync.WaitGroup

	// Start workers with their specific ranges
	for i := 1; i <= numWorkers; i++ {
		start, end := CalculateRange(i, numWorkers, maxNum)
		fmt.Printf("Worker %d assigned range: %d to %d\n", i, start, end)
		wg.Add(1)
		go Worker(i, start, end, results, progress, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Create a goroutine to collect and display progress
	workerProgress := make([]int, numWorkers+1) // +1 because worker ID start at 1
	totalWorkItems := make([]int, numWorkers+1)
	progressDone := make(chan bool)

	go func() {
		defer func() {
			progressDone <- true
		}()

		for update := range progress {
			workerProgress[update.WorkerID] = update.NumbersTested
			totalWorkItems[update.WorkerID] = update.TotalNumbers

			// Calculate overall progress
			totalProcessed := 0
			totalItems := 0
			for i := 1; i <= numWorkers; i++ {
				totalProcessed += workerProgress[i]
				totalItems += totalWorkItems[i]
			}

			if totalItems > 0 {
				percentage := float64(totalProcessed) * 100 / float64(totalItems)
				elapsed := time.Since(startTime)
				fmt.Printf("\rProgress: %.2f%% (%d/%d) - Elapsed: %s",
					percentage, totalProcessed, totalItems, elapsed.Round(time.Millisecond))
			}
		}
	}()

	// Collect prime numbers
	var primes []int
	for prime := range results {
		primes = append(primes, prime)
	}

	close(progress)

	// Wait for progress goroutine to finish
	<-progressDone

	// Print final results
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nFound %d prime numbers in %s\n", len(primes), elapsed.Round(time.Millisecond))

	// Sort primes (they will arrive out of order due to concurrent processing)
	sort.Ints(primes)

	// Print all prime numbers
	fmt.Println("\nAll prime numbers:")
	for i := 0; i < len(primes); i++ {
		fmt.Printf("%d ", primes[i])
		// Print a newline after every 10 numbers
		if (i+1)%10 == 0 {
			fmt.Println()
		}
	}

	// If the total count isn't a multiple of 10, add a final newline
	if len(primes)%10 != 0 {
		fmt.Println()
	}
}
