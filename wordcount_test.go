package main

import (
	"os"
	"path/filepath"
	"testing"
)

// //TESTS//
// Test case 1: Empty file
func TestEmptyFile(t *testing.T) {
	result := countWordsConcurrently("", 100)
	if result != 0 {
		t.Errorf("Expected 0 words for empty file, got %d", result)
	}
}

// Test case 2: Single word
func TestSingleWord(t *testing.T) {
	text := "Bonjour"
	result := countWordsConcurrently(text, 100)
	if result != 1 {
		t.Errorf("Expected 1 word for 'Bonjour', got %d", result)
	}
}

// Test case 3: Multiple lines with multiple words
func TestMultipleLinesMultipleWords(t *testing.T) {
	text := "Hello world\nThis is a test\nWith many words here"
	// Expected: Hello(1) world(2) This(3) is(4) a(5) test(6) With(7) many(8) words(9) here(10)
	expected := 10
	result := countWordsConcurrently(text, 100)
	if result != expected {
		t.Errorf("Expected %d words, got %d", expected, result)
	}
}

// Test case 4: Word boundary handling - segment cuts through word
func TestWordBoundaryHandling(t *testing.T) {
	text := "hello beautiful world"
	// With segment size 8, segments would be:
	// [0:8]   = "hello be" -> 1 word (hello, be is cut because it ends mid-word)
	// [8:16]  = "autiful " -> 1 words (counts word if starts mid-word)
	// [16:21] = "world"    -> 1 word
	// Total should be 3
	result := countWordsConcurrently(text, 8)
	if result != 3 {
		t.Errorf("Expected 3 words with segment size 8, got %d", result)
	}
}

// Test case 5: Single character segments
func TestSingleCharacterSegments(t *testing.T) {
	text := "hello world"
	// 11 characters, 2 words
	result := countWordsConcurrently(text, 1)
	if result != 2 {
		t.Errorf("Expected 2 words with char segments, got %d", result)
	}
}

// Test case 6: Large segment (entire text in one segment)
func TestLargeSegment(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	result := countWordsConcurrently(text, 1000)
	if result != 9 {
		t.Errorf("Expected 9 words with large segment, got %d", result)
	}
}

// Test case 7: Extra whitespace before and after text
func TestExtraWhiteSpaceAroundText(t *testing.T) {
	text := "  one two three four five    "
	result := countWordsConcurrently(text, 50)
	if result != 5 {
		t.Errorf("Expected 5 words, got %d", result)
	}
}

// Test case 8: Extra whitespace between words
func TestExtraWhitespace(t *testing.T) {
	text := "hello    world   test"
	// strings.Fields handles multiple spaces as single separator
	result := countWordsConcurrently(text, 100)
	if result != 3 {
		t.Errorf("Expected 3 words with extra whitespace, got %d", result)
	}
}

// //BENCHMARKS//
const benchmarkNormalTextSize = 10000   // 10k characters
const benchmarkLargeTextSize = 10000000 // 10 million characters

// Benchmark: Small segment size, normal text size (high concurrency)
func BenchmarkSmallSegments(b *testing.B) {
	text := generateTestText(benchmarkNormalTextSize)
	segmentSize := 10
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Medium segment size, normal text size
func BenchmarkMediumSegments(b *testing.B) {
	text := generateTestText(benchmarkNormalTextSize)
	segmentSize := 100
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Large segment size, normal text size (low concurrency)
func BenchmarkLargeSegments(b *testing.B) {
	text := generateTestText(benchmarkNormalTextSize)
	segmentSize := 1000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Very large segment size, normal text size (basically single-threaded)
func BenchmarkVeryLargeSegments(b *testing.B) {
	text := generateTestText(benchmarkNormalTextSize)
	segmentSize := 10000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Small segment size, large text size
func Benchmark10MBSmallSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 100
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Medium segment size, large text size
func Benchmark10MBMediumSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 1000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Large segment size, large text size
func Benchmark10MBLargeSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 10000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Very large segment size, large text size
func Benchmark10MBVeryLargeSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 100000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Mega segment size, large text size
func Benchmark10MBMegaSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 1000000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Benchmark: Gigantic segment size, large text size
func Benchmark10MBGiganticSegments(b *testing.B) {
	text := generateTestText(benchmarkLargeTextSize)
	segmentSize := 10000000
	b.ResetTimer()
	for b.Loop() {
		countWordsConcurrently(text, segmentSize)
	}
}

// Create test input files for manual testing
func init() {
	// Create test files in the same directory as the test
	testDir := "test_files"
	os.MkdirAll(testDir, 0755)

	// Empty file
	os.WriteFile(filepath.Join(testDir, "empty.txt"), []byte(""), 0644)

	// Single word file
	os.WriteFile(filepath.Join(testDir, "single_word.txt"), []byte("Bonjour"), 0644)

	// Multiple lines file
	multiContent := "Hello world\nThis is a test file\nWith multiple lines\nContaining many words\n"
	os.WriteFile(filepath.Join(testDir, "multiple_lines.txt"), []byte(multiContent), 0644)

	// Large file for performance testing
	largeContent := generateTestText(100000)
	os.WriteFile(filepath.Join(testDir, "large.txt"), []byte(largeContent), 0644)
}
