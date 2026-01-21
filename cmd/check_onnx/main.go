package main

import (
	"fmt"
	"os"
	"runtime"

	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	fmt.Println("Checking ONNX Runtime...")

	// Copy of logic from ai_model.go
	possiblePaths := []string{
		"/opt/homebrew/opt/onnxruntime/lib/libonnxruntime.dylib", // Apple Silicon Homebrew
		"/usr/local/opt/onnxruntime/lib/libonnxruntime.dylib",    // Intel Homebrew
		"/usr/local/lib/libonnxruntime.dylib",                    // Manual install
	}

	if runtime.GOOS == "linux" {
		possiblePaths = []string{
			"/usr/lib/libonnxruntime.so",
			"/usr/local/lib/libonnxruntime.so",
		}
	}

	var foundPath string
	for _, path := range possiblePaths {
		if _, e := os.Stat(path); e == nil {
			foundPath = path
			fmt.Printf("Found library at: %s\n", path)
			break
		}
	}

	if foundPath == "" {
		fmt.Println("Library path NOT found in listed locations.")
		os.Exit(1)
	}

	ort.SetSharedLibraryPath(foundPath)
	err := ort.InitializeEnvironment()
	if err != nil {
		fmt.Printf("FAIL: InitializeEnvironment returned error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS: ONNX Runtime initialized correctly.")
}
