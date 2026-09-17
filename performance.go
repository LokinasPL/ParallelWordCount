package main

import (
	"fmt"
	"strings"
	"time"
)

// Performance analysis function
// Run this separately to see how performance scales with segment size
func analyzePerformance() {
	testSizes := []struct {
		name        string
		fileSize    int
		segmentSize int
	}{
		{"10KB - Small segments", 10000, 10},
		{"10KB - Small-Medium segments", 10000, 50},
		{"10KB - Medium segments", 10000, 100},
		{"10KB - Medium-Large segments", 10000, 500},
		{"10KB - Large segments", 10000, 1000},
		{"10KB - Large-Very Large segments", 10000, 5000},
		{"10KB - Very large segments", 10000, 10000},

		{"100KB - Small segments", 100000, 100},
		{"100KB - Small-Medium segments", 100000, 500},
		{"100KB - Medium segments", 100000, 1000},
		{"100KB - Medium-Large segments", 100000, 5000},
		{"100KB - Large segments", 100000, 10000},
		{"100KB - Large-Very Large segments", 100000, 50000},
		{"100KB - Very large segments", 100000, 100000},

		{"1MB - Small segments", 1000000, 1000},
		{"1MB - Small-Medium segments", 1000000, 5000},
		{"1MB - Medium segments", 1000000, 10000},
		{"1MB - Medium-Large segments", 1000000, 50000},
		{"1MB - Large segments", 1000000, 100000},
		{"1MB - Large-Very Large segments", 1000000, 500000},
		{"1MB - Very large segments", 1000000, 1000000},
	}

	fmt.Println("=== PERFORMANCE ANALYSIS ===")
	fmt.Println("File Size | Segment Size | Num Segments | Avg Time (ns) | Words/ms")
	fmt.Println(strings.Repeat("-", 80))

	for _, test := range testSizes {
		text := generateTestText(test.fileSize)
		numSegments := (len(text) + test.segmentSize - 1) / test.segmentSize

		// Warm-up run
		countWordsConcurrently(text, test.segmentSize)

		// Time 10 iterations
		start := time.Now()
		for range 10 {
			countWordsConcurrently(text, test.segmentSize)
		}
		elapsed := time.Since(start)
		avgTimeNs := float64(elapsed.Nanoseconds()) / 10.0
		wordsPerMs := float64(len(strings.Fields(text))) / (avgTimeNs / 1e6)

		fmt.Printf("%8d | %12d | %12d | %13.3f | %10.2f\n",
			test.fileSize, test.segmentSize, numSegments, avgTimeNs, wordsPerMs)
	}
}
