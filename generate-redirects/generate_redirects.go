package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Path to the input file containing subdomains
const inputFilePathRelative = "../../fpf_frontend/src/core/signals/search.ts"

// Base directory to generate subdirectories and HTML files
const outputDirRelative = "../r/"

func main() {
	// Get the file path of the current file
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Unable to get the file path")
		return
	}
	dir := filepath.Dir(filePath)
	fmt.Println("Current file path:", dir)
	inputFilePath := filepath.Join(dir, inputFilePathRelative)
	fmt.Println("Input file path:", inputFilePath)
	outputDir := filepath.Join(dir, outputDirRelative)
	fmt.Println("Output directory:", outputDir)

	// Read the subdomains from the file
	subdomains, err := readSubdomains(inputFilePath)
	if err != nil {
		fmt.Printf("Error reading subdomains: %v\n", err)
		return
	}

	// Regex pattern to match subdomains ending with "-pornstar-lookalike"
	pattern := regexp.MustCompile(`([a-z\-]+)-pornstar-lookalike`)

	for _, subdomain := range subdomains {
		if matches := pattern.FindStringSubmatch(subdomain); matches != nil {
			prefix := matches[1] // The subdomain prefix (e.g., "subdomain1" from "subdomain1-article")
			// Create the subdomain folder and HTML file
			if err := createRedirectHTML(outputDir, prefix); err != nil {
				fmt.Printf("Error creating HTML for %s: %v\n", prefix, err)
			} else {
				fmt.Printf("Generated redirect for %s\n", prefix)
			}
		}
	}
}

// Reads subdomains from a file
func readSubdomains(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var subdomains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomains = append(subdomains, scanner.Text())
	}
	return subdomains, scanner.Err()
}

// Creates the subdomain folder and an index.html with the redirect
func createRedirectHTML(outputDir string, name string) error {
	// Define the target URL based on the subdomain
	targetURL := fmt.Sprintf("https://findpornface.com/article/%s-pornstar-lookalike", name)

	// Path for the subdomain folder
	folder := strings.ReplaceAll(name, "-", "")
	subdomainDir := filepath.Join(outputDir, folder)

	// Create the subdomain directory
	if err := os.MkdirAll(subdomainDir, os.ModePerm); err != nil {
		return err
	}

	// Path for the index.html file
	htmlFilePath := filepath.Join(subdomainDir, "index.html")

	// Write the redirect HTML content
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="refresh" content="0; url=%s">
    <title>Redirecting...</title>
  </head>
  <body>
    <p>If you are not redirected automatically, <a href="%s">click here</a>.</p>
  </body>
</html>`, targetURL, targetURL)

	return os.WriteFile(htmlFilePath, []byte(htmlContent), 0644)
}
