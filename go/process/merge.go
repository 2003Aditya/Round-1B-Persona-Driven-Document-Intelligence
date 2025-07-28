package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type PartialSection struct {
	Document     string  `json:"document"`
	SectionTitle string  `json:"section_title"`
	PageNumber   int     `json:"page_number"`
	Score        float64 `json:"score"`
	RefinedText  string  `json:"refined_text"`
}

type OutputJSON struct {
	Metadata struct {
		InputDocuments      []string `json:"input_documents"`
		Persona             string   `json:"persona"`
		JobToBeDone         string   `json:"job_to_be_done"`
		ProcessingTimestamp string   `json:"processing_timestamp"`
	} `json:"metadata"`
	ExtractedSections  []map[string]interface{} `json:"extracted_sections"`
	SubsectionAnalysis []map[string]interface{} `json:"subsection_analysis"`
}

func MergePartialOutputs(persona, job string, docs []DocumentInfo) {
	fmt.Println("🔧 Starting merge of partial outputs...")

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		sem      = make(chan struct{}, 4) // Limit to 4 concurrent readers
		partials []PartialSection
	)

	for _, doc := range docs {
		wg.Add(1)
		docCopy := doc // avoid race condition
		go func() {
			defer wg.Done()
			sem <- struct{}{} // acquire
			defer func() { <-sem }() // release

			baseName := strings.TrimSuffix(docCopy.Filename, ".pdf")
			filePath := filepath.Join("temp_output", fmt.Sprintf("%s_combined.partial.json", baseName))

			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("❌ Failed to read file %s: %v\n", filePath, err)
				return
			}
			fmt.Printf("📄 Reading partial output: %s\n", filePath)

			var sections []PartialSection
			if err := json.Unmarshal(data, &sections); err != nil {
				fmt.Printf("❌ Failed to unmarshal JSON from %s: %v\n", filePath, err)
				return
			}

			fmt.Printf("✅ Loaded %d sections from %s\n", len(sections), filePath)

			mu.Lock()
			partials = append(partials, sections...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Sort by score descending
	sort.Slice(partials, func(i, j int) bool {
		return partials[i].Score > partials[j].Score
	})

	// Deduplicate top 5 by section title and page number
	seen := make(map[string]bool)
	unique := []PartialSection{}
	for _, sec := range partials {
		key := fmt.Sprintf("%s|%d", sec.SectionTitle, sec.PageNumber)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, sec)
		}
		if len(unique) >= 5 {
			break
		}
	}

	// Prepare output
	out := OutputJSON{}
	for _, d := range docs {
		out.Metadata.InputDocuments = append(out.Metadata.InputDocuments, d.Filename)
	}
	out.Metadata.Persona = persona
	out.Metadata.JobToBeDone = job
	out.Metadata.ProcessingTimestamp = time.Now().Format(time.RFC3339)

	for i, sec := range unique {
		out.ExtractedSections = append(out.ExtractedSections, map[string]interface{}{
			"document":        sec.Document,
			"section_title":   sec.SectionTitle,
			"page_number":     sec.PageNumber,
			"importance_rank": i + 1,
		})
		out.SubsectionAnalysis = append(out.SubsectionAnalysis, map[string]interface{}{
			"document":     sec.Document,
			"page_number":  sec.PageNumber,
			"refined_text": sec.RefinedText,
		})
	}

	outputPath := "output/challenge1b_output.json"
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("❌ Failed to create output file: %v\n", err)
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Printf("❌ Failed to write JSON output: %v\n", err)
		return
	}

	fmt.Printf("✅ Output successfully written to %s\n", outputPath)
}

