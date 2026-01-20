package game

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestMultiUserQueueLoad simulates multiple concurrent users hitting the global service
func TestMultiUserQueueLoad(t *testing.T) {
	numUsers := 3
	onnxPath := "../../ml/checkpoints/snake_policy.onnx"

	// Start the global service
	err := StartInferenceService(onnxPath)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(numUsers)

	fmt.Printf("🔥 Starting QUEUE-BASED Load Test with %d concurrent users...\n", numUsers)
	start := time.Now()

	for i := 0; i < numUsers; i++ {
		go func(id int) {
			defer wg.Done()

			input := make([]float64, 6*25*25)
			for j := 0; j < 100; j++ {
				// Global Predict function handles queueing automatically
				Predict(input)
			}
			fmt.Printf("✅ User %d: Finished 100 queued inferences.\n", id)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	totalInferences := numUsers * 100
	avgPerInference := elapsed / time.Duration(totalInferences)

	fmt.Printf("\n--- Global Queue Service Result ---\n")
	fmt.Printf("Total Concurrent Users: %d\n", numUsers)
	fmt.Printf("Total Inferences: %d\n", totalInferences)
	fmt.Printf("Total Time Elapsed: %v\n", elapsed)
	fmt.Printf("System Throughput: %.2f inferences/sec\n", float64(totalInferences)/elapsed.Seconds())
	fmt.Printf("Average Latency (End-to-end): %v\n", avgPerInference)
	fmt.Printf("-----------------------------------\n")
}
