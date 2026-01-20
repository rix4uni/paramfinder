package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
	"strings"

	"github.com/spf13/pflag"
	"github.com/rix4uni/paramfinder/banner"
)

func main() {
	// Define command-line flags
	numRoutines := pflag.Int("concurrency", 50, "number of concurrent goroutines")
	timeout := pflag.Int("timeout", 30, "HTTP request timeout duration (in seconds)")
	outputFileFlag := pflag.String("output", "", "output file path")
	silent := pflag.Bool("silent", false, "silent mode.")
	version := pflag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := pflag.Bool("verbose", false, "enable verbose mode")

	// Parse the command-line flags
	pflag.Parse()

	if *version {
		banner.PrintBanner()
		banner.PrintVersion()
		os.Exit(0)
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Create a multi-writer for output
	var outputWriter io.Writer = os.Stdout
	if *outputFileFlag != "" {
		output, err := os.Create(*outputFileFlag)
		if err != nil {
			fmt.Println("Error opening output file:", err)
			os.Exit(1)
		}
		defer output.Close()
		outputWriter = io.MultiWriter(os.Stdout, output)
	}

	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to send URLs to be processed
	urlChan := make(chan string)

	// Create an HTTP client with the specified timeout and insecure setting enabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Duration(*timeout) * time.Second,
		Transport: tr,
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
				re := regexp.MustCompile(`<input[^>]*>|<textarea[^>]*>`)
				inputTags := re.FindAllString(string(body), -1)

				// Transform URL and print it
				transformedURL := notransformURLWithParams(url, inputTags)
				if transformedURL != url { // Check if transformation resulted in a different URL
					fmt.Fprintln(outputWriter, transformedURL)
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

// notransformURLWithParams appends query parameters to the URL based on input tags
func notransformURLWithParams(baseURL string, inputTags []string) string {
	// Create an ordered map to keep track of the parameters and their values
	params := make([]string, 0)
	seen := make(map[string]bool)

	for _, tag := range inputTags {
		re := regexp.MustCompile(`name="([^"]+)"`)
		names := re.FindAllStringSubmatch(tag, -1)
		for _, name := range names {
			paramName := name[1]
			if !seen[paramName] {
				params = append(params, fmt.Sprintf("%s=rix4uni", paramName))
				seen[paramName] = true
			}
		}
	}

	queryString := strings.Join(params, "&")

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	parsedURL.RawQuery = queryString
	return parsedURL.String()
}
