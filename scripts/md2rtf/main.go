package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	pdf "github.com/stephenafamo/goldmark-pdf"
	rtfdoc "github.com/therox/rtf-doc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	gohtml "golang.org/x/net/html"
)

var rootDir = filepath.Join("..", "..")

// Escape RTF special characters
func escapeRTF(input string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"{", "\\{",
		"}", "\\}",
	)
	return replacer.Replace(input)
}

const baseURL = "https://raw.githubusercontent.com/BrycensRanch/Rokon/refs/heads/master/"

// Fetch image bytes from a URL or local path
func fetchImage(src string, rootDir string) ([]byte, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		req, err := http.NewRequest("GET", src, nil)
		if err != nil {
			return nil, err
		}

		// Set the Accept-Encoding header to 'identity'
		req.Header.Set("Accept-Encoding", "identity;q=0")
		resp, err := http.DefaultClient.Do(req)
		var body io.Reader = resp.Body
		if resp.Header.Get("Content-Encoding") == "br" {
			// Decode the Brotli response
			body = brotli.NewReader(resp.Body)
		} else {
			body = resp.Body
		}

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(body)
	} else {
		// Handle local image paths
		localImagePath := filepath.Join(rootDir, src)
		return ioutil.ReadFile(localImagePath)
	}
}

// Decompress gzipped data
func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

// Convert image bytes to PNG if necessary and add it to the RTF document
func addImageToRTF(doc *rtfdoc.Document, imgBytes []byte) error {
	imgType := http.DetectContentType(imgBytes)
	var img image.Image
	var err error

	switch imgType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imgBytes))
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imgBytes))
	case "image/svg+xml":
		imgBytes, err = convertSVGToPNG(imgBytes)
		if err != nil {
			return err
		}
		img, err = png.Decode(bytes.NewReader(imgBytes))
		imgType = "image/png"
	case "text/plain; charset=utf-8":
		imgBytes, err = convertSVGToPNG(imgBytes)
		if err != nil {
			return err
		}
		img, err = png.Decode(bytes.NewReader(imgBytes))
		imgType = "image/png"
	case "application/octet-stream": // snap store bs the web server is a fucking liar
		uncompressedBytes, err := decompressGzip(imgBytes)
		if err != nil {
			return err
		}
		// Check if it's SVG after decompression
		if strings.Contains(http.DetectContentType(uncompressedBytes), "image/svg") {
			imgBytes, err = convertSVGToPNG(uncompressedBytes)
			if err != nil {
				return err
			}
			img, err = png.Decode(bytes.NewReader(imgBytes))
		} else {
			return fmt.Errorf("unsupported image format: %s", imgType)
		}

	default:
		return fmt.Errorf("unsupported image format: %s", imgType)
	}

	if err != nil {
		log.Printf("addImageToRTF ran into error:")
		return err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	var format string
	if strings.Contains(imgType, "png") {
		format = "png"
	} else {
		format = "jpeg"
	}

	doc.AddParagraph().AddPicture(buf.Bytes(), format) // Adjust width and height as necessary
	return nil
}

// Convert SVG bytes to PNG bytes using ImageMagick
func convertSVGToPNG(svgBytes []byte) ([]byte, error) {
	tempSVGFile, err := ioutil.TempFile("", "image-*.svg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempSVGFile.Name()) // Clean up the temp file

	if _, err := tempSVGFile.Write(svgBytes); err != nil {
		return nil, err
	}
	tempSVGFile.Close()

	tempPNGFile, err := ioutil.TempFile("", "image-*.png")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempPNGFile.Name()) // Clean up the temp file

	cmd := exec.Command("magick", tempSVGFile.Name(), tempPNGFile.Name())
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(tempPNGFile.Name())
}

// Convert Markdown to HTML
func convertMarkdownToHTML(md string) (string, error) {
	var buf bytes.Buffer
	mdParser := goldmark.New(
		goldmark.WithRendererOptions(html.WithUnsafe(), html.WithHardWraps()),
	)

	if err := mdParser.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Convert HTML to RTF using rtf-doc
func convertHTMLToRTF(htmlContent string) (string, error) {
	doc := rtfdoc.NewDocument()
	doc.SetOrientation(rtfdoc.OrientationPortrait)
	doc.SetFormat(rtfdoc.FormatA4)

	node, err := gohtml.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var f func(*gohtml.Node)
	f = func(n *gohtml.Node) {
		if n.Type == gohtml.ElementNode {
			align := rtfdoc.AlignLeft // Default alignment
			// Check for align attribute
			for _, attr := range n.Attr {
				if attr.Key == "align" {
					switch attr.Val {
					case "center":
						align = rtfdoc.AlignCenter
					case "right":
						align = rtfdoc.AlignRight
					case "justify":
						align = rtfdoc.AlignJustify
					}
				}
			}
			switch n.Data {
			case "h1":
				p := doc.AddParagraph()
				p.AddText(escapeRTF(getText(n)), 44, rtfdoc.FontCourierNew, rtfdoc.ColorBlack)
				p.SetAlign(align)
			case "h2":
				p := doc.AddParagraph()
				p.AddText(escapeRTF(getText(n)), 32, rtfdoc.FontCourierNew, rtfdoc.ColorBlue)
				p.SetAlign(align)
			case "h3":
				p := doc.AddParagraph()
				p.AddText(escapeRTF(getText(n)), 24, rtfdoc.FontCourierNew, rtfdoc.ColorOlive)
				p.SetAlign(align)
			case "p":
				p := doc.AddParagraph()
				p.AddText(escapeRTF(getText(n)), 12, rtfdoc.FontCourierNew, rtfdoc.ColorBlack)
				p.SetAlign(align)

			case "b":
				p := doc.AddParagraph()
				txt := p.AddText(escapeRTF(getText(n)), 12, rtfdoc.FontCourierNew, rtfdoc.ColorBlack)
				txt.SetBold()
				p.SetAlign(align)
			case "i":
				p := doc.AddParagraph()
				txt := p.AddText(escapeRTF(getText(n)), 12, rtfdoc.FontCourierNew, rtfdoc.ColorBlack)
				txt.SetItalic()
				p.SetAlign(align)
			case "a":
				p := doc.AddParagraph()
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						txt := p.AddText(escapeRTF(getText(n)), 12, rtfdoc.FontCourierNew, rtfdoc.ColorBlue)
						txt.SetBold()
						p.SetAlign(align)
					}
				}
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						log.Printf(attr.Val)
						imgBytes, err := fetchImage(attr.Val, rootDir)
						if err != nil {
							fmt.Println("Error fetching image:", err)
							continue
						}
						if err := addImageToRTF(doc, imgBytes); err != nil {
							fmt.Println("Error adding image to RTF:", err)
						}
					}
				}
			case "br":
				doc.AddParagraph().AddNewLine()
			case "li":
				p := doc.AddParagraph()
				p.AddText(escapeRTF(getText(n)), 12, rtfdoc.FontCourierNew, rtfdoc.ColorBlack)
				p.SetAlign(rtfdoc.AlignLeft)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)
	return string(doc.Export()), nil
}

func getText(n *gohtml.Node) string {
	var buf strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == gohtml.TextNode {
			buf.WriteString(c.Data)
		} else if c.Type == gohtml.ElementNode {
			buf.WriteString(getText(c))
		}
	}
	return buf.String()
}

func main() {
	if os.Args[1] == "md2pdf" {
		md2pdf()
		return
	}
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			mdFilePath := filepath.Join(rootDir, file.Name())
			rtfFilePath := strings.TrimSuffix(mdFilePath, ".md") + ".rtf"

			mdContent, err := ioutil.ReadFile(mdFilePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			htmlContent, err := convertMarkdownToHTML(string(mdContent))
			if err != nil {
				fmt.Println("Error converting Markdown to HTML:", err)
				continue
			}

			rtfContent, err := convertHTMLToRTF(htmlContent)
			if err != nil {
				fmt.Println("Error converting HTML to RTF:", err)
				continue
			}

			err = ioutil.WriteFile(rtfFilePath, []byte(rtfContent), 0o644)
			if err != nil {
				fmt.Println("Error writing RTF file:", err)
			}

			fmt.Printf("Converted %s to %s\n", mdFilePath, rtfFilePath)
		}
	}
}

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

func md2pdf() {
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
