package process

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
)

type DocumentInfo struct {
    Filename string
    Title    string
}

func ProcessPDF(filename string) {
    inputPath := filepath.Join("input", filename)
    chunkOutDir := "temp_output"
    os.MkdirAll(chunkOutDir, 0755)

    // Step 1: Count total pages
    cmd := exec.Command("python3", "python/extractor/extract.py", "--count", inputPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("❌ Failed to count pages in %s: %v", filename, err)
        fmt.Println("Python error:\n" + string(output))
        return
    }
    totalPages, err := strconv.Atoi(strings.TrimSpace(string(output)))
    if err != nil {
        log.Printf("❌ Invalid page count output: %s", output)
        return
    }
    chunkSize := (totalPages + 4) / 5

    // Step 2: Extract chunks concurrently
    var chunkWg sync.WaitGroup
    for i := 0; i < 5; i++ {
        start := i * chunkSize
        end := min(start+chunkSize, totalPages)
        if start >= end {
            continue
        }

        outFile := fmt.Sprintf("%s_chunk_%d.json", strings.TrimSuffix(filename, ".pdf"), i)
        outPath := filepath.Join(chunkOutDir, outFile)

        chunkWg.Add(1)
        go func(start, end int, outPath string) {
            defer chunkWg.Done()
            cmd := exec.Command("python3", "python/extractor/extract.py", inputPath, fmt.Sprint(start), fmt.Sprint(end), outPath)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            if err := cmd.Run(); err != nil {
                log.Printf("❌ Chunk extraction failed (%d-%d): %v", start, end, err)
            } else {
                log.Printf("✅ Saved to %s", outPath)
            }
        }(start, end, outPath)
    }
    chunkWg.Wait()

    // Step 3: Combine all chunk files into one
    combinedFile := fmt.Sprintf("%s_combined.json", strings.TrimSuffix(filename, ".pdf"))
    combinedPath := filepath.Join(chunkOutDir, combinedFile)

    combineCmd := exec.Command("python3", "python/extractor/extract.py", "--combine", strings.TrimSuffix(filename, ".pdf"), combinedPath)
    combineCmd.Stdout = os.Stdout
    combineCmd.Stderr = os.Stderr
    if err := combineCmd.Run(); err != nil {
        log.Printf("❌ Failed to combine chunks for %s: %v", filename, err)
        return
    }
    log.Printf("✅ Combined chunks → %s", combinedPath)
}

// ✅ Call this *once* after all ProcessPDF calls
func RunScoringBatch(persona, job string) {
    scorerCmd := exec.Command(
        "python3", "python/analyzer/scorer.py",
        "--batch_dir", "temp_output",
        "--persona", persona,
        "--job", job,
    )
    scorerCmd.Stdout = os.Stdout
    scorerCmd.Stderr = os.Stderr
    if err := scorerCmd.Run(); err != nil {
        log.Printf("❌ Scoring batch failed: %v", err)
    } else {
        log.Printf("✅ Scoring completed in parallel.")
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

