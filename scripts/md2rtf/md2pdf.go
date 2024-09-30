package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func convertMarkdownToPDF(mdContent []byte, outputPath string) error {
	// Create a new goldmark Markdown parser with the PDF renderer
	md := goldmark.New(
		goldmark.WithRenderer(
			pdf.New(
				pdf.WithContext(context.Background()),
				// pdf.WithImageFS(os.DirFS("../..")),
				pdf.WithHeadingFont(pdf.GetTextFont("IBM Plex Serif", pdf.FontLora)),
				pdf.WithBodyFont(pdf.GetTextFont("Helvetica", pdf.FontHelvetica)),
				pdf.WithCodeFont(pdf.GetCodeFont("Fire Code", pdf.FontFiraCode)),
			),
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	// Convert the Markdown content to PDF and write it to the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create PDF file: %w", err)
	}
	defer file.Close()

	if err := md.Convert(mdContent, file); err != nil {
		return fmt.Errorf("could not convert markdown to PDF: %w", err)
	}

	return nil
}

func main() {
	rootDir := filepath.Join("..", "..")

	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			mdFilePath := filepath.Join(rootDir, file.Name())
			pdfFilePath := strings.TrimSuffix(mdFilePath, ".md") + ".pdf"

			mdContent, err := ioutil.ReadFile(mdFilePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			if err := convertMarkdownToPDF(mdContent, pdfFilePath); err != nil {
				fmt.Println("Error converting to PDF:", err)
				continue
			}

			fmt.Printf("Converted %s to %s\n", mdFilePath, pdfFilePath)
		}
	}
}
