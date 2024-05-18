package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

func main() {
	// Define command-line flags
	numRoutines := flag.Int("c", 20, "number of concurrent goroutines")
	timeout := flag.Int("timeout", 30, "HTTP request timeout duration (in seconds)")
	verbose := flag.Bool("v", false, "enable verbose mode")
	outputFile := flag.String("o", "", "output file path")

	// Parse the command-line flags
	flag.Parse()

	// Create a multi-writer for output
	var outputWriter io.Writer = os.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Println("Error opening output file:", err)
			os.Exit(1)
		}
		defer file.Close()
		outputWriter = io.MultiWriter(os.Stdout, file)
	}

	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to send URLs to be processed
	urlChan := make(chan string)

	// Create an HTTP client with the specified timeout
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	// Start the goroutines
	for i := 0; i < *numRoutines; i++ {
		wg.Add(1)
		go func() {
			// Decrement the wait group counter when the goroutine finishes
			defer wg.Done()

			// Process URLs from the channel
			for url := range urlChan {
				// Make an HTTP GET request to the URL
				resp, err := client.Get(url)
				if err != nil {
					if *verbose {
						fmt.Fprintln(outputWriter, err)
					}
					continue
				}
				defer resp.Body.Close()

				// Read the response body into a string
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					if *verbose {
						fmt.Fprintln(outputWriter, err)
					}
					continue
				}

				// Use a regular expression to find all input tags in the body
				re := regexp.MustCompile(`<input[^>]*>`)
				inputTags := re.FindAllString(string(body), -1)

				// Print the URL and input tags if verbose mode is enabled or there are input tags
				if *verbose || len(inputTags) > 0 {
					fmt.Fprintln(outputWriter, "URL:", url)
					for _, tag := range inputTags {
						fmt.Fprintln(outputWriter, tag)
					}
					fmt.Fprintln(outputWriter)
				}
			}
		}()
	}

	// Loop through each line (URL) in standard input and send it to the channel
	for scanner.Scan() {
		url := scanner.Text()
		urlChan <- url
	}

	// Close the channel to indicate that there are no more URLs to process
	close(urlChan)

	// Wait for all goroutines to finish
	wg.Wait()

	// Check for errors while scanning standard input
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
