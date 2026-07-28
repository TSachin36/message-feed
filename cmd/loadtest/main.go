package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	target      = "http://localhost:8080/messages?user=alice"
	requests    = 1000
	concurrency = 20
)

func main() {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	jobs := make(chan int)

	var wg sync.WaitGroup

	var success atomic.Int64
	var failed atomic.Int64

	start := time.Now()

	for i := 0; i < concurrency; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {

				response, err := client.Get(target)
				if err != nil {
					failed.Add(1)
					continue
				}

				_, err = io.Copy(
					io.Discard,
					response.Body,
				)

				response.Body.Close()

				if err != nil {
					failed.Add(1)
					continue
				}

				if response.StatusCode >= 200 &&
					response.StatusCode < 300 {

					success.Add(1)

				} else {
					failed.Add(1)
				}
			}
		}()
	}

	for i := 0; i < requests; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()

	duration := time.Since(start)

	fmt.Println("\nLoad Test Results")
	fmt.Println("-----------------")

	fmt.Printf("Requests:    %d\n", requests)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Successful:  %d\n", success.Load())
	fmt.Printf("Failed:      %d\n", failed.Load())
	fmt.Printf("Duration:    %v\n", duration)

	fmt.Printf(
		"Requests/sec: %.2f\n",
		float64(requests)/duration.Seconds(),
	)

	fmt.Printf(
		"Average time: %v\n",
		duration/time.Duration(requests),
	)
}
