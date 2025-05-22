package main

import (
	"sync"
	"testing"
)

func TestIsPrime(t *testing.T) {
	nonPrimes := []int{-1, 0, 1, 4, 6, 9, 15, 25}
	for _, n := range nonPrimes {
		if IsPrime(n) {
			t.Errorf("%d incorrectly identified as prime", n)
		}
	}

	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23}
	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("%d incorrectly identified as non-prime", n)
		}
	}
}

func TestCalculateRange(t *testing.T) {
	testCases := []struct {
		worker int
		start  int
		end    int
	}{
		{1, 1, 25},
		{2, 26, 50},
		{3, 51, 75},
		{4, 76, 100},
	}

	for _, tc := range testCases {
		start, end := CalculateRange(tc.worker, 4, 100)
		if start != tc.start || end != tc.end {
			t.Errorf("Worker %d: got range (%d, %d), expected (%d, %d)",
				tc.worker, start, end, tc.start, tc.end)
		}
	}
}

func TestWorker(t *testing.T) {
	done := make(chan bool)
	results := make(chan int, 5)
	progress := make(chan ProgressUpdate, 5)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		Worker(1, 10, 20, results, progress, &wg)
		close(results)
		close(progress)
		done <- true
	}()

	var primes []int
	for prime := range results {
		primes = append(primes, prime)
	}

	<-done 

	expected := []int{11, 13, 17, 19}
	
	if len(primes) != len(expected) {
		t.Errorf("Found %d primes, expected %d", len(primes), len(expected))
		return
	}

	primesFound := make(map[int]bool)
	for _, p := range primes {
		primesFound[p] = true
	}

	for _, p := range expected {
		if !primesFound[p] {
			t.Errorf("Prime %d not found", p)
		}
	}
}
