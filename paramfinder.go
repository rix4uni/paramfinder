package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Set the number of goroutines to 8
	numRoutines := 20

	// Create a channel to send URLs to be processed
	urlChan := make(chan string)

	// Start the goroutines
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			// Decrement the wait group counter when the goroutine finishes
			defer wg.Done()

			// Process URLs from the channel
			for url := range urlChan {
				// Make an HTTP GET request to the URL
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				defer resp.Body.Close()

				// Read the response body into a string
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// Use a regular expression to find all input tags in the body
				re := regexp.MustCompile(`<input[^>]*>`)
				inputTags := re.FindAllString(string(body), -1)

				// Print the URL and input tags
				fmt.Println(url)
				for _, tag := range inputTags {
					fmt.Println(tag)
				}
				fmt.Println()
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
