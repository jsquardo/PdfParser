package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run pdf_replace.go <search_text> <replace_text> <pdf_file1> [<pdf_file2> ...]")
		os.Exit(1)
	}

	searchText := os.Args[1]
	replaceText := os.Args[2]
	pdfFiles := os.Args[3:]

	// Set the unipdf log level to error to reduce noise.
	common.SetLogger(common.ConsoleLogger{LogLevel: common.LogError})

	for _, pdfFile := range pdfFiles {
		fmt.Println("Processing:", pdfFile)

		// Read the PDF file into a PDF model object.
		pdf, err := model.NewPdfReaderFromFile(pdfFile)
		if err != nil {
			fmt.Println("Error reading PDF:", err)
			continue
		}

		// Create a new PDF creator object with the same configuration as the input PDF.
		c := creator.New()

		// Loop through each page of the PDF.
		for i := 1; i <= pdf.GetNumPages(); i++ {
			fmt.Printf("Processing page %d\n", i)

			// Get the page object and content stream.
			page := pdf.GetPage(i)
			contentStream, err := page.GetAllContentStreams()
			if err != nil {
				fmt.Println("Error getting content stream:", err)
				continue
			}

			// Convert the content stream to a string and replace the search text with the replace text.
			contentStreamStr := contentStream.String()
			newContentStreamStr := strings.ReplaceAll(contentStreamStr, searchText, replaceText)

			// Add the modified content stream to the new PDF creator object.
			c.NewPage()
			c.WriteText(newContentStreamStr)
		}

		// Save the modified PDF to a new file.
		outputFilename := pdfFile + ".new"
		err = c.WriteToFile(outputFilename)
		if err != nil {
			fmt.Println("Error saving PDF:", err)
			continue
		}
	}

	fmt.Println("Done!")
}
