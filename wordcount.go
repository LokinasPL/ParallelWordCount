package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

/*
// if testing the performance only, uncomment this version of main and comment the normal main
func main() {
	analyzePerformance()
}*/

func main() {
	//check that no arg is missing
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file> <segment-size>\n", os.Args[0])
		os.Exit(1)
	}

	//get the file path arg
	filePath := os.Args[1]
	//get the segment size arg
	segmentSize, err := strconv.Atoi(os.Args[2])
	//check that the segment size is a positive integer and that there was no error in parsing it
	if err != nil || segmentSize <= 0 {
		fmt.Fprintf(os.Stderr, "segment-size must be a positive integer\n")
		os.Exit(1)
	}

	//read the file content
	content, err := os.ReadFile(filePath)
	//check that there was no error in reading the file
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		os.Exit(1)
	}

	//call concurrent word counter and print the result
	totalWords := countWordsConcurrently(string(content), segmentSize)
	fmt.Printf("Total words: %d\n", totalWords)
}

// counts the number of words in the given text using concurrent goroutines to process segments of the text defined by the segment size
func countWordsConcurrently(text string, segmentSize int) int {
	//convert the text to a slice of runes to handle Unicode characters properly
	runes := []rune(text)
	//get the length of the runes slice
	length := len(runes)
	//if there is no character in the text, return 0 words counted
	if length == 0 {
		return 0
	}

	//create a channel to collect the word counts from each segment
	results := make(chan int)
	//create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	//for each segment of the text
	for start := 0; start < length; start += segmentSize {
		//calculate the position of the end of the segment, making sure it doesn't go beyond the length of the text
		end := start + segmentSize
		if end > length {
			end = length
		}

		//launch a goroutine to count the words in the current segment
		wg.Add(1)
		//launch goroutine with the start and end positions of the segment as arguments
		go func(start, end int) {
			//defer the wg.Done call so that this waitgroup will end even if there is a panic in the goroutine
			defer wg.Done()
			//call countSegmentWords and send the result to the results channel
			results <- countSegmentWords(runes, start, end)
		}(start, end)
	}

	//once all the segements have been processed, call a last goroutine
	go func() {
		//we wait for all goroutines in the wait group to finish
		wg.Wait()
		//then close the results channel
		close(results)
	}()

	//initialize total word count to 0
	total := 0
	//sum all the counts sent to the result channel
	for count := range results {
		total += count
	}

	//return the total word count
	return total
}

// count the number of words in a segment of the text defined by the start and end positions in the runes slice
func countSegmentWords(runes []rune, start, end int) int {
	//initialize the segment as a string from the runes slice fromg the start to the end positions
	segment := string(runes[start:end])
	//initialize the word count as the numbre of words in the segment string
	words := len(strings.Fields(segment))

	//check if the segment string ends mid-word
	endsMidWord := end < len(runes) && !unicode.IsSpace(runes[end-1]) && !unicode.IsSpace(runes[end])

	//if the segment ends in the middle of a word, we don't count that word in this segment, it will be counted in the next segment
	if endsMidWord {
		words--
	}
	//if by some freak case the word count is negative, we set it to 0 , since a negative value is impossible
	if words < 0 {
		words = 0
	}
	//return the word count for this segment
	return words
}

// Helper function to generate test text, used in tests, benchmarks, and performance analysis
func generateTestText(size int) string {
	words := []string{"hello", "world", "test", "concurrent", "goroutine", "channel", "segment", "performance", "benchmark", "analysis"}
	text := ""
	wordCount := 0
	charCount := 0

	for charCount < size {
		text += words[wordCount%len(words)] + " "
		charCount += len(words[wordCount%len(words)]) + 1
		wordCount++
	}

	return text[:min(len(text), size)]
}
