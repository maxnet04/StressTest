package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

const (
	timeoutDuration = 3 * time.Second
	statusTimeout   = 408
)

var cmd = &cobra.Command{

	Use:   "stress-test-cli",
	Short: "A cli to do stress test",
	Long:  `A cli to do a stress test. Setting up the --url <target> to work`,
	Run:   executeStressTest,
}

func executeStressTest(cmd *cobra.Command, args []string) {

	start := time.Now()

	url, requests, concurrency, err := parseFlags(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Performing load test on %s with %d requests and %d concurrent cells.\n", url, requests, concurrency)

	statusCode := performLoadTest(url, requests, concurrency)

	elapsed := time.Since(start)
	generateReport(elapsed, requests, &statusCode)

}

func parseFlags(cmd *cobra.Command) (string, int, int, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil || url == "" {
		return "", 0, 0, errors.New("url is required")
	}

	requests, err := cmd.Flags().GetInt("requests")
	if err != nil {
		return "", 0, 0, err
	}
	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return "", 0, 0, err
	}
	return url, requests, concurrency, nil
}

func performLoadTest(url string, requests, concurrency int) sync.Map {

	var wg sync.WaitGroup
	statusCodes := sync.Map{}
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
			defer cancel()

			res, err := makeRequest(&ctx, url)
			if err != nil {
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					summary(&statusCodes, statusTimeout)
				} else {
					updateSummary(&statusCodes, statusTimeout)
				}
				return
			}
			summary(&statusCodes, res.StatusCode)
		}()
	}
	wg.Wait()
	return statusCodes
}

func makeRequest(ctx *context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(*ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

func generateReport(executionTime time.Duration, numberOfRequests int, statusCodes *sync.Map) {
	fmt.Printf("==============================================\n")
	fmt.Printf("Total execution time: %s\n", executionTime)
	fmt.Printf("Total number of request made: %d\n", numberOfRequests)
	fmt.Printf("\nSummary: \n")
	statusCodes.Range(func(key, value interface{}) bool {
		count := value.(int32)
		percentage := (float64(count) * 100 / float64(numberOfRequests))
		fmt.Printf("Status code %d | Count  %d (%.2f%%)\n", key, count, percentage)
		return true
	})
}

func updateSummary(statusCodes *sync.Map, key int) {
	existingCount, _ := statusCodes.LoadOrStore(key, int32(0))
	statusCodes.Store(key, existingCount.(int32)+1)
}

func Execure() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().StringP("url", "u", "", "Define the target url")
	cmd.Flags().IntP("requests", "r", 100, "Define the number of requests")
	cmd.Flags().IntP("concurrency", "c", 1, "Define the number of concurrent requests")
}
