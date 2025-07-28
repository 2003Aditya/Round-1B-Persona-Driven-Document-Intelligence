import argparse
import os
import sys
import json
import fitz  # PyMuPDF
import spacy
import traceback
import time
from functools import lru_cache
from sklearn.metrics.pairwise import cosine_similarity
from concurrent.futures import ProcessPoolExecutor, as_completed

parser = argparse.ArgumentParser()
parser.add_argument("--batch_dir", type=str)
parser.add_argument("--persona", type=str)
parser.add_argument("--job", type=str)
parser.add_argument("args", nargs="*")
args = parser.parse_args()

try:
    nlp = spacy.load("en_core_web_md")
except:
    print("⚠️ Falling back to en_core_web_sm")
    nlp = spacy.load("en_core_web_sm")

@lru_cache(maxsize=2048)
def get_vector(text: str):
    return nlp(text).vector.reshape(1, -1)

def score_file(task):
    filename, persona, job, input_path, output_path = task
    try:
        if os.path.exists(output_path):
            print(f"⚠️ Skipping {output_path}, already exists.")
            return output_path

        print(f"📑 Scoring sections from: {input_path}")
        start = time.time()

        query_vec = get_vector(f"{persona} wants to {job}")

        with open(input_path) as f:
            sections = json.load(f)

        doc = fitz.open(filename)
        scored = []

        for section in sections:
            title = section["text"]
            page_number = section["page"]
            section_vec = get_vector(title)
            section_score = cosine_similarity(query_vec, section_vec).item()

            try:
                page_text = doc.load_page(page_number).get_text()
            except:
                page_text = ""

            sentences = [sent.text.strip() for sent in nlp(page_text).sents if sent.text.strip()]
            scored_sents = []
            for idx, sent in enumerate(sentences):
                sent_vec = get_vector(sent)
                sent_score = cosine_similarity(query_vec, sent_vec).item()
                scored_sents.append((idx, sent_score, sent))

            relevant = [(i, s) for i, score, s in scored_sents if score > 0.6]
            relevant.sort(key=lambda x: x[0])
            refined_text = " ".join(s for _, s in relevant[:10]) if relevant else title

            scored.append({
                "document": os.path.basename(filename),
                "section_title": title,
                "page_number": page_number,
                "score": section_score,
                "refined_text": refined_text
            })

        with open(output_path, "w") as f:
            json.dump(scored, f, indent=2)

        print(f"✅ Scored {len(scored)} sections → {output_path} in {round(time.time() - start, 2)}s")
        return output_path

    except Exception as e:
        print(f"❌ Error while scoring {input_path}: {e}")
        traceback.print_exc()
        return None

def run_batch_scoring(batch_dir, persona, job, max_workers=8):
    files = [f for f in os.listdir(batch_dir) if f.endswith("_combined.json")]
    print(f"🚀 Starting scoring for {len(files)} files using {max_workers} workers...")

    tasks = []
    for file in files:
        base = file.replace("_combined.json", "")
        pdf_path = os.path.join("input", base + ".pdf")
        in_path = os.path.join(batch_dir, file)
        out_path = os.path.join(batch_dir, base + "_combined.partial.json")
        tasks.append((pdf_path, persona, job, in_path, out_path))

    with ProcessPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(score_file, task) for task in tasks]
        for future in as_completed(futures):
            _ = future.result()

if args.batch_dir:
    run_batch_scoring(args.batch_dir, args.persona, args.job)

elif len(args.args) == 5:
    score_file(tuple(args.args))

else:
    print("❌ Invalid arguments.")
    sys.exit(1)

