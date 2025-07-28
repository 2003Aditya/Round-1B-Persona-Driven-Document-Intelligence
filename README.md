
# 🧠 Round 1B: Persona-Driven Document Intelligence

---

## 🚀 Overview

Build an *intelligent document analyzer* that extracts and ranks the most relevant sections from a PDF set based on:
- 🧑 A specific **persona**
- 🎯 Their **job to be done**

> 📌 Theme: "Connect What Matters — For the User Who Matters"

---

## 📥 Input

- Folder: `input/`
- File: `challenge1b_input.json`

```json
{
  "persona": "Investment Analyst",
  "job_to_be_done": "Analyze revenue trends, R&D investment, market strategy",
  "input_documents": ["report1.pdf", "report2.pdf", "report3.pdf"]
}
```

---

## 📤 Output

- File: `output/challenge1b_output.json`

```json
{
  "metadata": {
    "input_documents": ["report1.pdf", "report2.pdf"],
    "persona": "Investment Analyst",
    "job_to_be_done": "Analyze revenue trends, R&D investment, market strategy",
    "processing_timestamp": "2025-07-28T16:32:45Z"
  },
  "extracted_sections": [
    {
      "document": "report1.pdf",
      "section_title": "Revenue Analysis",
      "page_number": 5,
      "importance_rank": 1
    }
  ],
  "subsection_analysis": [
    {
      "document": "report1.pdf",
      "page_number": 5,
      "refined_text": "The company reported a YoY growth of..."
    }
  ]
}
```

✅ Flat, clean, persona-focused, and sorted by importance.

---

## 🧠 Methodology

1. 📚 **Parse PDFs** using PyMuPDF
2. 🧠 **Embed persona + job description** using SpaCy (en_core_web_md)
3. 📄 **Extract headings & sections** from all PDFs
4. 🧮 **Compute cosine similarity** between persona-job vector and each section
5. 🏆 **Rank top results** based on relevance score
6. ✂️ **Extract snippets** from ranked heading

---

## 🧪 Sample Use Cases

| Persona               | Job-To-Be-Done                                                 |
|-----------------------|----------------------------------------------------------------|
| PhD Student           | Literature review on GNN benchmarks                            |
| Investment Analyst    | Analyze revenue, R&D, strategy in annual reports               |
| Chemistry Student     | Study reaction kinetics concepts from textbook PDFs            |

---

## 🐳 Run via Docker (Offline / amd64)

```bash
# Prepare
mkdir -p input output temp_output
cp your_pdfs.pdf input/
cp challenge1b_input.json input/

# Build
docker build -t doc-intelligence .

# Run
docker run --rm \
  -v $(pwd)/input:/app/input \
  -v $(pwd)/output:/app/output \
  -v $(pwd)/temp_output:/app/temp_output \
  doc-intelligence --input input/challenge1b_input.json
```

---

## 📁 Project Structure

```
.
├── Dockerfile
├── go/
│   ├── cmd/main.go
│   └── process/
├── python/
│   ├── extractor/extract.py
│   ├── analyzer/scorer.py
│   └── requirements.txt
├── input/
├── output/
├── temp_output/
└── README.md
```

---

## ⚙️ Tech Stack

- Python 3.10
- PyMuPDF for parsing
- SpaCy (en_core_web_md) for semantic similarity
- Go for orchestration and concurrency

---

## ✅ Hackathon Checklist

| Constraint                  | Status     |
|----------------------------|------------|
| Runtime ≤ 60s              | ✅ Pass     |
| Model < 1GB                | ✅ Pass     |
| No Internet                | ✅ Pass     |
| Platform (amd64)           | ✅ Docker OK|
| Persona-Driven Ranking     | ✅ Yes      |
| Clean JSON Output          | ✅ Pass     |

---

## 📬 Contact

Made with 💡 for Adobe India Hackathon – Round 1B
GitHub: [https://github.com/2003Aditya]

