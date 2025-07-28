package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/2003Aditya/process"
)

type Input struct {
	Persona     struct{ Role string } `json:"persona"`
	JobToBeDone struct{ Task string } `json:"job_to_be_done"`
	Documents   []struct {
		Filename string `json:"filename"`
		Title    string `json:"title"`
	} `json:"documents"`
}

func main() {
	start := time.Now()

	// ✅ Allow dynamic input file via flag
	var inputFilePath string
	flag.StringVar(&inputFilePath, "input", "input/challenge1b_input.json", "Path to the input JSON file")
	flag.Parse()

	inputPath := filepath.Clean(inputFilePath)
	inputFile, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to read input JSON: %v", err)
	}

	var input Input
	if err := json.Unmarshal(inputFile, &input); err != nil {
		log.Fatalf("❌ Invalid input JSON: %v", err)
	}

	fmt.Println("🚀 Starting Challenge 1B Engine")
	fmt.Printf("📌 Persona: %s\n🛠️  Task: %s\n\n", input.Persona.Role, input.JobToBeDone.Task)

	var globalWg sync.WaitGroup
	for _, doc := range input.Documents {
		docCopy := doc
		globalWg.Add(1)
		go func(doc process.DocumentInfo) {
			defer globalWg.Done()
			process.ProcessPDF(doc.Filename)
		}(process.DocumentInfo{
			Filename: docCopy.Filename,
			Title:    docCopy.Title,
		})
	}
	globalWg.Wait()

	// Merge + score once
	var documents []process.DocumentInfo
	for _, doc := range input.Documents {
		documents = append(documents, process.DocumentInfo{
			Filename: doc.Filename,
			Title:    doc.Title,
		})
	}

	process.RunScoringBatch(input.Persona.Role, input.JobToBeDone.Task)
	process.MergePartialOutputs(input.Persona.Role, input.JobToBeDone.Task, documents)

	fmt.Println("✅ challenge1b_output.json created successfully.")
	fmt.Printf("⏱️  Finished in %v\n", time.Since(start))
}

