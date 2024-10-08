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
	"math/rand"

	"github.com/spf13/pflag"
)

// prints the version message
const version = "0.0.2"

func printVersion() {
	fmt.Printf("Current paramfinder version %s\n", version)
}

// Prints the banner
func printBanner() {
	banner := `
                                           ____ _             __           
    ____   ____ _ _____ ____ _ ____ ___   / __/(_)____   ____/ /___   _____
   / __ \ / __  // ___// __  // __  __ \ / /_ / // __ \ / __  // _ \ / ___/
  / /_/ // /_/ // /   / /_/ // / / / / // __// // / / // /_/ //  __// /    
 / .___/ \__,_//_/    \__,_//_/ /_/ /_//_/  /_//_/ /_/ \__,_/ \___//_/     
/_/                                                                    `
fmt.Printf("%s\n%80s\n\n", banner, "Current paramfinder version "+version)

}

// Generate a random string of lowercase letters of the specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	// Define command-line flags
	numRoutines := pflag.IntP("concurrency", "c", 50, "number of concurrent goroutines")
	timeout := pflag.IntP("timeout", "t", 10, "HTTP request timeout duration (in seconds)")
	outputFileFlag := pflag.StringP("output", "o", "", "output file path")
	appendOutputFlag := pflag.StringP("append", "a", "", "File to append the output instead of overwriting.")
	insecure := pflag.BoolP("insecure", "i", false, "allow insecure server connections when using SSL")
	notransformURL := pflag.BoolP("no-turl", "n", false, "Do not print transform URL with extracted parameters")
	onlyHidden := pflag.Bool("only-hidden", false, "print only hidden input tags")
	silent := pflag.BoolP("silent", "s", false, "silent mode.")
	version := pflag.BoolP("version", "V", false, "Print the version of the tool and exit.")
	verbose := pflag.BoolP("verbose", "v", false, "enable verbose mode")

	// Parse the command-line flags
	pflag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silent flag is provided
	if !*silent {
		printBanner()
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
	} else if *appendOutputFlag != "" {
		output, err := os.OpenFile(*appendOutputFlag, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening output file for appending:", err)
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

	// Create an HTTP client with the specified timeout and optional insecure setting
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
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

				// Filter for only hidden input tags if the flag is set
				if *onlyHidden {
					var hiddenTags []string
					for _, tag := range inputTags {
						if strings.Contains(tag, `type="hidden"`) || strings.Contains(tag, `type='hidden'`) {
							hiddenTags = append(hiddenTags, tag)
						}
					}
					inputTags = hiddenTags
				}

				// Print the URL and input tags if verbose mode is enabled or there are input tags
				if *verbose || len(inputTags) > 0 {
					fmt.Fprintln(outputWriter, "URL:", url)
					for _, tag := range inputTags {
						fmt.Fprintln(outputWriter, tag)
					}
					fmt.Fprintln(outputWriter)
				}

				// Transform URL if -no-turl flag is not set
				if !*notransformURL {
					transformedURL := notransformURLWithParams(url, inputTags, *onlyHidden)
					if transformedURL != url { // Check if transformation resulted in a different URL
						fmt.Fprintln(outputWriter, "TRANSFORM_URL:", transformedURL)
					}
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
func notransformURLWithParams(baseURL string, inputTags []string, onlyHidden bool) string {
	// Create an ordered map to keep track of the parameters and their values
	params := make([]string, 0)
	seen := make(map[string]bool)

	for _, tag := range inputTags {
		re := regexp.MustCompile(`name="([^"]+)"`)
		names := re.FindAllStringSubmatch(tag, -1)
		for _, name := range names {
			paramName := name[1]
			if !seen[paramName] {
				randomString := generateRandomString(7)
				params = append(params, fmt.Sprintf("%s=%s", paramName, randomString))
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
