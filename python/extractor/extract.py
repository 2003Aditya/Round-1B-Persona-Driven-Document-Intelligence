import fitz  # PyMuPDF
import json
import os
import sys

def extract_headings(pdf_path, start_page=0, end_page=None):
    doc = fitz.open(pdf_path)
    title = doc.metadata.get("title", "") or doc.load_page(0).get_text().split('\n')[0]
    outline = []

    if end_page is None:
        end_page = len(doc)

    for page_num in range(start_page, min(end_page, len(doc))):
        page = doc.load_page(page_num)
        blocks = page.get_text("dict")["blocks"]
        for block in blocks:
            for line in block.get("lines", []):
                spans = line.get("spans", [])
                if not spans:
                    continue
                text = " ".join(span["text"].strip() for span in spans).strip()
                if not text or len(text.split()) > 12 or len(text) < 5:
                    continue
                font_sizes = [round(span["size"]) for span in spans]
                if max(font_sizes) >= 12:
                    outline.append({
                        "level": "H1",
                        "text": text,
                        "page": page_num
                    })
    doc.close()
    return {"title": title.strip(), "outline": outline}

def combine_chunks(filename_prefix, output_path):
    combined = []
    for i in range(5):
        chunk_path = f"temp_output/{filename_prefix}_chunk_{i}.json"
        if not os.path.exists(chunk_path):
            continue
        with open(chunk_path) as f:
            data = json.load(f)
            combined.extend(data)
    with open(output_path, "w") as f:
        json.dump(combined, f, indent=2)
    print(f"✅ Combined all chunks → {output_path}")



def main():
    args = sys.argv[1:]

    if "--count" in args:
        if len(args) < 2:
            print("Missing PDF path for --count")
            sys.exit(1)
        doc = fitz.open(args[1])
        print(len(doc))
        return

    elif "--combine" in args:
        if len(args) < 3:
            print("Missing arguments for --combine <filename_prefix> <output_path>")
            sys.exit(1)
        filename_prefix, output_path = args[1], args[2]
        combine_chunks(filename_prefix, output_path)
        return

    elif len(args) == 4:
        pdf_path, start_page, end_page, output_path = args
        start_page, end_page = int(start_page), int(end_page)
        result = extract_headings(pdf_path, start_page, end_page)
        with open(output_path, "w") as f:
            json.dump(result["outline"], f, indent=2)
        print(f"Saved to {output_path}")
        return

    else:
        print("Usage:")
        print("  python extract.py <pdf> <start_page> <end_page> <output_path>")
        print("  python extract.py --count <pdf>")
        print("  python extract.py --combine <prefix> <output_path>")

if __name__ == "__main__":
    main()

