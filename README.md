# AfriLearn Curriculum API

> **The open digital backbone & curriculum data layer for African educational technology, school portals, edtech apps, and AI Tutors.**
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB UTME · NUC CCMAS Degrees · NBTE Polytechnics · Open Data (MIT License) · AI Tutor LLM Prompts · Interactive Swagger Docs**

---

## 🌍 Open Data & Community Commitment

All raw curriculum datasets in `data/curricula/` are **100% open-source under the MIT License**. Anyone can clone, use, extend, or contribute datasets to improve educational data access across Africa.

- **Datasets Location**: [`data/curricula/`](./data/curricula/) (56 JSON files covering 14 institutions & exam boards)
- **Contribution Guide**: [`CONTRIBUTING.md`](./CONTRIBUTING.md) — schema rules, validation CLI, and PR workflow
- **Agentic Evaluation Report**: [`questionstest.md`](./questionstest.md) — 6-level parallel comparative test (Baseline LLM vs AfriLearn Activated)
- **Technical Roadmap**: [`ROADMAP.md`](./ROADMAP.md) — 6-phase engineering plan (embeddings, progress tracking, adaptive pathways)

---

## 🌟 Capabilities Overview

| Feature | Status | Engineering Detail |
|---|---|---|
| 56 structured curriculum datasets | ✅ Production | BECE, WAEC, JAMB, NUC, NBTE + 14 universities (UNILAG, FUTO, UNN, etc.) |
| 3-layer full-text search | ✅ Production | PostgreSQL GIN indexes across topics, subtopics, and learning objectives |
| LLM system prompt engine | ✅ Production | Bloom's Taxonomy breakdown + domain rules (Law, Med, Eng, Nursing, CS) |
| Cross-board curriculum matcher | ✅ Production | `GET /api/v1/curriculum/match/:topic` — queries all 22 boards simultaneously |
| Learning pathway engine | ✅ Production | BECE → WAEC → JAMB → University progression ordering |
| Prerequisite graph engine | ✅ Production | `GET /api/v1/curriculum/prerequisites/:board/:subject/:topic` |
| Natural language query brain | ✅ Production | Intent parser with grade-level and topic inference fallbacks |
| Admin operations dashboard | ✅ Production | `/admin` — real-time cache stats, purge, and background re-ingestion |
| Interactive AI Playground | ✅ Production | `/playground` — test all endpoints with live AI tutor chat (Groq/Gemini) |
| Interactive Swagger API docs | ✅ Production | `/docs` — OpenAPI spec and curl generator |
| In-memory caching layer | ✅ Production | Thread-safe, per-prefix hit/miss metrics (`curr:`, `prompt:`, `embeddings:`) |
| PGVector HNSW index | ✅ Schema-ready | Provisioned column (`topics.embedding`) ready for ML vector upserts |
| RAG text chunking engine | ✅ Production | `GET /api/v1/curriculum/:board/:subject/embeddings` — 200–800 token semantic blocks |
| **Real semantic vector search** | 🔄 Phase 2 | Scheduled: OpenAI / Gemini / Ollama embedding pipeline (see ROADMAP.md) |

> [!NOTE]
> **Search Engine Implementation Note**: The search endpoints (`GET /api/v1/search` and `POST /api/v1/search/vector`) perform **ranked PostgreSQL full-text search (FTS)** across topics, subtopics, and objectives. The `pgvector` HNSW column is schema-ready for Phase 2 when OpenAI/Gemini embeddings are upserted into `topics.embedding`.

---

## 🚀 Quick Start

### 1. Clone & Set Up
```bash
git clone https://github.com/HarperKollins/AfriLearn.git
cd AfriLearn
go mod tidy
cp .env.example .env
```

### 2. Configure Database
AfriLearn uses PostgreSQL with the `pgvector` extension. Neon.tech (free tier), Supabase, or local Postgres all work. Set `DB_URL` in `.env`:
```ini
DB_URL="postgresql://user:password@localhost:5432/afrilearn?sslmode=disable"
```

### 3. Deploy Schema & Ingest Curricula
```bash
# Deploys schema, tables, GIN full-text indexes, and PGVector HNSW index
go run cmd/migrate/main.go

# Ingests all 56 curriculum files (~45 seconds, idempotent — safe to re-run)
go run cmd/seeder/main.go

# Optional: validate all datasets for schema correctness first
go run cmd/seeder/main.go --validate-only
```

### 4. Run Automated Test Suite
```bash
go test -v ./...
```

### 5. Start the API Server
```bash
go run cmd/api/main.go
```

| URL | Purpose |
|---|---|
| `http://localhost:8080/` | Developer portal |
| `http://localhost:8080/admin` | Admin operations dashboard |
| `http://localhost:8080/playground` | Interactive AI Playground |
| `http://localhost:8080/docs` | Swagger / OpenAPI interactive docs |
| `http://localhost:8080/health` | Health check |

---

## 🔑 Authentication

All API endpoints (except health and key generation) require an `X-API-Key` header.

**Demo key** (rate-limited, for testing): `afr_live_demo_9f8e2b7a`

**Generate a key**:
```bash
curl -X POST http://localhost:8080/api/v1/keys/generate \
  -H "Content-Type: application/json" \
  -d '{"name": "My App", "email": "dev@example.com"}'
```

---

## 📚 API Endpoint Reference & Payload Samples

### 1. Get Full Curriculum Tree
Returns a board's complete curriculum tree: topics → subtopics → learning objectives.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/unilag/law
```

### 2. Get AI Tutor System Prompt
Returns an LLM system prompt with Bloom's Taxonomy breakdown, domain-specific rules (Law, Medicine, Engineering, CS, etc.), difficulty progression, and token count.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/unilag/law/llm-prompt
```

**Response Payload Sample**:
```json
{
  "success": true,
  "data": {
    "exam_board": "University of Lagos",
    "subject": "Law (LL.B.)",
    "system_prompt": "You are an expert AI Tutor specialized in the official University of Lagos Law curriculum...",
    "subject_specific_rules": [
      "Apply the IRAC method for problem questions: Issue → Rule → Application → Conclusion.",
      "Always cite the full case name, court, and year (e.g., Donoghue v Stevenson [1932] AC 562 (HL)).",
      "Reference the Nigerian Constitution 1999 and CAMA 2020 by section number."
    ],
    "blooms_taxonomy_breakdown": {
      "remember": 12, "understand": 34, "apply": 28, "analyze": 15, "evaluate": 8, "create": 3
    },
    "estimated_token_count": 4200,
    "suggested_chunking_note": "~4200 tokens — fits in most 8K context models. Use full_context_window directly."
  }
}
```

### 3. Get RAG Embedding Chunks
Returns curriculum content pre-chunked (one chunk per topic module, 200–800 tokens) ready to embed with OpenAI, Gemini, or Ollama into Pinecone, Qdrant, ChromaDB, or PGVector.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/waec/physics/embeddings
```

**Response Payload Sample**:
```json
{
  "success": true,
  "data": {
    "board": "WAEC",
    "subject": "Physics",
    "total_chunks": 6,
    "total_token_estimate": 320,
    "chunking_strategy": "one-chunk-per-topic-module (200–800 tokens per chunk)",
    "chunks": [
      {
        "chunk_id": "waec_physics_module_01",
        "module_title": "Interaction of Matter, Force, and Motion",
        "text_content": "CURRICULUM: West African Examinations Council — Physics\nLEVEL: SS1-SS3 | BOARD: WAEC\n\nMODULE 1: Interaction of Matter, Force, and Motion\nDIFFICULTY: INTERMEDIATE\n...",
        "token_estimate": 55,
        "metadata": {
          "board_full_name": "West African Examinations Council",
          "module_index": 1,
          "difficulty": "intermediate",
          "action_verbs": ["calculate", "define", "explain"]
        }
      }
    ],
    "integration_guide": {
      "openai": "client.embeddings.create(model='text-embedding-3-small', input=chunk['text_content'])",
      "google_gemini": "genai.embed_content(model='models/text-embedding-004', content=chunk['text_content'])",
      "ollama": "ollama.embeddings(model='nomic-embed-text', prompt=chunk['text_content'])"
    }
  }
}
```

### 4. Search Across All Curricula (3-Layer Deep FTS)
Searches topics, subtopics, AND learning objectives simultaneously with ranked results.
```bash
# Basic search
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/search?q=photosynthesis"

# Filtered search (by board and subject) with pagination
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/search?q=newton+laws&board=waec&subject=physics&limit=10&offset=0"
```

### 5. Vector Search Endpoint (API Contract)
```bash
curl -X POST -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  -H "Content-Type: application/json" \
  -d '{"query": "statutory interpretation rules in legal methods", "limit": 5}' \
  http://localhost:8080/api/v1/search/vector
```

---

## ⚡ Performance Benchmarking & Load Testing

AfriLearn includes a built-in CLI load test utility to measure throughput and latency distribution (p50, p90, p99) under concurrent worker load.

```bash
# Start server in terminal 1
go run cmd/api/main.go

# Run load benchmark in terminal 2 (10 concurrent workers, 200 requests per endpoint)
go run cmd/loadtest/main.go -url http://localhost:8080 -c 10 -n 200
```

**Sample Benchmark Output**:
```text
🚀 AfriLearn API Load Test & Latency Benchmark
   Base URL:     http://localhost:8080
   Concurrency:  10 workers
   Total Reqs:   200 per test
──────────────────────────────────────────────────────
✅ API Health check passed

⚡ Benchmarking: GET Curriculum (WAEC Math)
   Success:    200/200 (Errors: 0)
   Throughput: 852.14 req/sec (Total Time: 234ms)
   Latency:    Min: 1.1ms | p50: 8.4ms | p90: 18.2ms | p99: 34.1ms | Max: 42.0ms

⚡ Benchmarking: GET Deep FTS Search ('quadratic')
   Success:    200/200 (Errors: 0)
   Throughput: 142.30 req/sec (Total Time: 1.40s)
   Latency:    Min: 12.4ms | p50: 58.1ms | p90: 112.4ms | p99: 185.0ms | Max: 210.3ms
```

---

## 🤖 Agentic LLM Evaluation & System Activation Benchmark

AfriLearn includes an automated comparative evaluation tool in `cmd/agentic_eval/main.go` that runs parallel tests comparing **Mode A (Unassisted General LLM Baseline)** against **Mode B (AfriLearn System Activated)** across 6 curriculum levels (BECE, WAEC, JAMB, UNILAG Law, FUTO Engineering, NUC Computer Science).

```bash
# Run the agentic evaluation suite against the live database
go run cmd/agentic_eval/main.go
```

The tool queries live system directives from the database, runs comparative evaluation, and generates the complete audit report in [`questionstest.md`](./questionstest.md).

---

## 🗃️ Supported Boards & Institutions (22 total)

| Slug | Institution / Board | Level |
|---|---|---|
| `bece` | NERDC / Junior WAEC | Junior Secondary (JSS1–3) |
| `waec` | West African Examinations Council | Senior Secondary (SS1–3) |
| `neco` | National Examinations Council | Senior Secondary (SS1–3) |
| `jamb` | Joint Admissions & Matriculation Board | UTME Entrance |
| `nuc` | National Universities Commission (CCMAS Standard) | Tertiary Degree |
| `nbte` | National Board for Technical Education | ND / HND |
| `unilag` | University of Lagos | Tertiary Degree |
| `unn` | University of Nigeria, Nsukka | Tertiary Degree |
| `unec` | University of Nigeria, Enugu Campus | Tertiary Degree |
| `ebsu` | Ebonyi State University | Tertiary Degree |
| `funai` | Federal University, Ndufu-Alike Ikwo | Tertiary Degree |
| `futo` | Federal University of Technology, Owerri | Tertiary Degree |
| `yabatech` | Yaba College of Technology | ND / HND |
| `imt` | Institute of Management & Technology, Enugu | ND / HND |

---

## 🛠️ Admin Dashboard (`/admin`)

Access `http://localhost:8080/admin` for operational metrics:
- **Cache Statistics**: Hit ratio, total entries, per-prefix stats (`curr:`, `prompt:`, `embeddings:`)
- **Cache Control**: Instantly flush cache entries
- **Background Re-Ingestion**: Re-seed all 56 datasets from disk without server downtime
- **Dataset Validation**: Run schema compliance checks across all JSON files

API endpoints:
```bash
GET  /api/v1/admin/cache/stats     # Real-time cache metrics
POST /api/v1/admin/cache/purge     # Flush memory cache
POST /api/v1/admin/reingest        # Trigger dataset re-ingestion
GET  /api/v1/admin/validate        # Run schema validation checks
```

---

## ⚙️ Environment Variables

```ini
DB_URL=postgresql://user:password@host:5432/afrilearn?sslmode=require
PORT=8080
APP_ENV=development   # or: production
```

See [`.env.example`](./.env.example) for reference settings.

---

## ⚖️ License

The AfriLearn Curriculum API and all raw curriculum datasets in `data/curricula/` are distributed under the **MIT License**.
