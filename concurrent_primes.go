package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func IsPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsConcurrentPrime(n int, pchannel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return
		}
	}
	pchannel <- n
}

func Primes(n int) []int {
	primes := make([]int, 0)
	for i := 2; i < n; i++ {
		if IsPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func ConcurrentPrimes(n int) []int {
	primes := make([]int, 0)
	pchannel := make(chan int, n)
	wg := &sync.WaitGroup{}
	for i := 2; i < n; i++ {
		wg.Add(1)
		go IsConcurrentPrime(i, pchannel, wg)
	}

	go func(wg *sync.WaitGroup, pchannel chan int) {
		wg.Wait()
		close(pchannel)
	}(wg, pchannel)

	for p := range pchannel {
		primes = append(primes, p)
	}
	return primes
}

func main() {
	var limit = 10000000
	start := time.Now()
	Primes(limit)
	fmt.Printf("Sequential primes calculation %v\n", time.Since(start))

	start = time.Now()
	ConcurrentPrimes(limit)
	fmt.Printf("Concurrent prime calculation %v\n", time.Since(start))
}
