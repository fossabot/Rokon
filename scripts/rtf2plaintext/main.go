package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/huantt/plaintext-extractor"
)

func main() {
	// Get the root directory (two directories up from current directory)
	rootDir := filepath.Join("..", "..")
	outputDir := filepath.Join(rootDir, "windows") // Update this to your actual path, assuming the script is two directories deep

	// Use filepath.Glob to find all *.md files in the root directory
	markdownFiles, err := filepath.Glob(filepath.Join(rootDir, "*.md"))
	if err != nil {
		log.Fatalf("Failed to list markdown files: %v", err)
	}

	// Create a new Markdown extractor
	extractor := plaintext.NewMarkdownExtractor()

	for _, mdFile := range markdownFiles {
		// Read the content of each Markdown file
		content, err := os.ReadFile(mdFile)
		if err != nil {
			log.Printf("Failed to read file %s: %v", mdFile, err)
			continue
		}

		// Convert the Markdown content to plain text
		outputPtr, err := extractor.PlainText(string(content))
		if err != nil {
			log.Printf("Failed to extract plain text from %s: %v", mdFile, err)
			continue
		}

		// Dereference the output pointer to get the actual string value
		output := *outputPtr

		// Define the output .txt file path
		txtFileName := filepath.Base(mdFile[:len(mdFile)-3]) + ".txt" // Get the base file name and replace .md with .txt
		txtFilePath := filepath.Join(outputDir, txtFileName)          // Combine the output directory and file name

		// Write the plain text to the new .txt file
		if err := ioutil.WriteFile(txtFilePath, []byte(output), 0o644); err != nil {
			log.Printf("Failed to write to file %s: %v", txtFilePath, err)
			continue
		}

		log.Printf("Converted %s to %s\n", mdFile, txtFilePath)
	}
}
