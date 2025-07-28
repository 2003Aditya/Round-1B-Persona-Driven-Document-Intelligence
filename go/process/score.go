package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ScoreCombinedFile runs scorer.py for a given combined JSON file
func ScoreCombinedFile(filePath, persona, job string) error {
	base := filepath.Base(filePath)
	title := strings.TrimSuffix(base, "_combined.json")
	pdfFile := title + ".pdf"
	outputPath := strings.Replace(filePath, "_combined.json", "_combined.partial.json", 1)

	cmd := exec.Command(
		"python3", "scorer.py",
		pdfFile,
		persona,
		job,
		filePath,
		outputPath,
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scoring failed for %s: %v", filePath, err)
	}
	return nil
}

// ScoreAllCombinedFiles detects and scores all *_combined.json files in temp_output concurrently
func ScoreAllCombinedFiles(persona, job string) {
	dir := "../temp_output"
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("❌ Failed to read directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // limit concurrency

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_combined.json") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())

		wg.Add(1)
		sem <- struct{}{}

		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Printf("🚀 Starting scoring process...\n📑 Scoring sections from: %s\n", path)

			if err := ScoreCombinedFile(path, persona, job); err != nil {
				fmt.Printf("❌ Failed to score: %v\n", err)
			} else {
				fmt.Printf("✅ Scored → %s\n", strings.Replace(path, "_combined.json", "_combined.partial.json", 1))
			}
		}(filePath)
	}

	wg.Wait()
}

