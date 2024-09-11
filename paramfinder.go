package main

import (
	"bufio"
	"crypto/tls"
	"flag"
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
)

// prints the version message
const version = "0.0.1"

func printVersion() {
	fmt.Printf("Current paramfinder version %s\n", version)
}

// Prints the banner
func printBanner() {
	banner := `
 ____   __    ____    __    __  __  ____  ____  _  _  ____  ____  ____ 
(  _ \ /__\  (  _ \  /__\  (  \/  )( ___)(_  _)( \( )(  _ \( ___)(  _ \
 )___//(__)\  )   / /(__)\  )    (  )__)  _)(_  )  (  )(_) ))__)  )   /
(__) (__)(__)(_)\_)(__)(__)(_/\/\_)(__)  (____)(_)\_)(____/(____)(_)\_)`
fmt.Printf("%s\n%60s\n\n", banner, "Current paramfinder version "+version)

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
	numRoutines := flag.Int("c", 50, "number of concurrent goroutines")
	timeout := flag.Int("timeout", 10, "HTTP request timeout duration (in seconds)")
	outputFileFlag := flag.String("o", "", "output file path")
	appendOutputFlag := flag.String("ao", "", "File to append the output instead of overwriting.")
	insecure := flag.Bool("insecure", false, "allow insecure server connections when using SSL")
	transformURL := flag.Bool("turl", false, "transform URL with extracted parameters")
	silent := flag.Bool("silent", false, "silent mode.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := flag.Bool("verbose", false, "enable verbose mode")

	// Parse the command-line flags
	flag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silnet flag is provided
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

				// Print the URL and input tags if verbose mode is enabled or there are input tags
				if *verbose || len(inputTags) > 0 {
					fmt.Fprintln(outputWriter, "URL:", url)
					for _, tag := range inputTags {
						fmt.Fprintln(outputWriter, tag)
					}
					fmt.Fprintln(outputWriter)
				}

				// Transform URL if -turl flag is set
				if *transformURL {
					transformedURL := transformURLWithParams(url, inputTags)
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

// transformURLWithParams appends query parameters to the URL based on input tags
func transformURLWithParams(baseURL string, tags []string) string {
	// Create an ordered map to keep track of the parameters and their values
	params := make([]string, 0)
	seen := make(map[string]bool)

	for _, tag := range tags {
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
