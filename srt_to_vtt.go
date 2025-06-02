package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// for "sourceFiles", set the path to your SRT files directory which you downloaded from the customer
	sourceFiles := "~/Downloads/srtfiles"

	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run srt_to_vtt.go %s\n", sourceFiles)
		return
	}

	inputDir := os.Args[1]
	outputDir := filepath.Join(inputDir, "vtt-files")

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.srt files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".srt") {
			return nil
		}

		fmt.Println("Converting:", info.Name())
		return convertSRTtoVTT(path, filepath.Join(outputDir, strings.TrimSuffix(info.Name(), ".srt")+".vtt"))
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
	} else {
		fmt.Println("All .srt files converted and saved to", outputDir)
	}
}

func convertSRTtoVTT(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", outputPath, err)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	writer.WriteString("WEBVTT\n\n")

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			line = strings.ReplaceAll(line, ",", ".")
		}
		writer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", inputPath, err)
	}

	return writer.Flush()
}
