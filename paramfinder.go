package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Loop through each line (URL) in standard input
	for scanner.Scan() {
		url := scanner.Text()

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

	// Check for errors while scanning standard input
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
