package game

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkAIInference measures the speed of the queue-based ONNX service
func BenchmarkAIInference(b *testing.B) {
	// 1. Setup: Start Global Service
	onnxPath := "../../ml/checkpoints/snake_policy.onnx"
	err := StartInferenceService(onnxPath)
	if err != nil {
		b.Fatalf("Failed to start service: %v", err)
	}

	// 2. Prepare mock input (6x25x25)
	input := make([]float64, 6*25*25)

	// 3. Warm up
	Predict(input)

	// 4. Run Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Predict(input)
	}
}

// TestSingleRunReport provides a readable single-pass execution report
func TestSingleRunReport(t *testing.T) {
	onnxPath := "../../ml/checkpoints/snake_policy.onnx"
	err := StartInferenceService(onnxPath)
	if err != nil {
		t.Skip("Service failed to start")
	}

	input := make([]float64, 6*25*25)

	// Warm up to ensure model is in memory
	Predict(input)

	start := time.Now()
	logits := Predict(input)
	elapsed := time.Since(start)

	fmt.Printf("\n--- Global AI Worker Performance Report ---\n")
	fmt.Printf("Model: %s\n", onnxPath)
	fmt.Printf("Single Inference (Round-trip): %v\n", elapsed)
	fmt.Printf("Throughput Cap: %.2f FPS\n", 1.0/elapsed.Seconds())
	fmt.Printf("Logic Output: %v\n", logits)
	fmt.Printf("-------------------------------------------\n")
}
