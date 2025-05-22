package main

import (
	"math"
	"sync"
)

const (
	DEFAULT_MAX_NUMBER     = 100 // Default upper limit for finding primes
	DEFAULT_NUMBER_WORKERS = 4   // Default number of worker goroutines
)

// ProgressUpdate represents a progress update from a worker
type ProgressUpdate struct {
	WorkerID      int
	TotalNumbers  int
	NumbersTested int
}

// IsPrime checks if a number is prime using an optimized algorithm
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 5; i <= sqrt; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// CalculateRange determines the range of numbers a specific worker should process
func CalculateRange(workerID, numWorkers, maxNumber int) (int, int) {
	rangeSize := maxNumber / numWorkers
	start := 1 + (workerID-1)*rangeSize
	end := workerID * rangeSize

	// Handle remainder for last worker
	if workerID == numWorkers {
		end = maxNumber
	}

	return start, end
}

// Worker processes a range of numbers, finds primes, and reports progress
func Worker(id, start, end int, results chan<- int, progress chan<- ProgressUpdate, wg *sync.WaitGroup) {
	defer wg.Done()

	// Calculate total numbers in this range (inclusive of both start and end)
	totalNumbers := end - start + 1
	numbersTested := 0

	for num := start; num <= end; num++ {
		if IsPrime(num) {
			results <- num
		}

		numbersTested++

		// Report progress periodically (every 1000 numbers or at completion)
		if numbersTested%1000 == 0 || numbersTested == totalNumbers {
			progress <- ProgressUpdate{
				WorkerID:      id,
				TotalNumbers:  totalNumbers,
				NumbersTested: numbersTested,
			}
		}
	}
}
