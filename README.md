# AfriLearn Curriculum API

> **The open digital backbone & curriculum data layer for African educational technology, school portals, edtech apps, and AI Tutors.**
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB UTME · NUC CCMAS Degrees · NBTE Polytechnics · Open Data (MIT License) · AI Tutor LLM Prompts · Interactive Swagger Docs**

---

## 🌍 Open Data & Community Commitment

All raw curriculum datasets in `data/curricula/` are **100% open-source under the MIT License**. Anyone can clone, use, extend, or contribute datasets to improve educational data access across Africa.

- **Datasets Location**: [`data/curricula/`](./data/curricula/) — 56 JSON files covering 14 institutions & exam boards
- **Contribution Guide**: [`CONTRIBUTING.md`](./CONTRIBUTING.md) — schema rules, validation CLI, and PR workflow
- **Benchmark Report**: [`BENCHMARK_REPORT.md`](./BENCHMARK_REPORT.md) — Claude Sonnet 4.6 baseline vs AfriLearn System (16 scenarios, +90.4% gain)
- **Full Evaluation Suite**: [`questionstest.md`](./questionstest.md) — 1,147-line parallel A/B test across 6 curriculum levels
- **Technical Roadmap**: [`ROADMAP.md`](./ROADMAP.md) — 6-phase engineering plan (embeddings, progress tracking, adaptive pathways)

---

## 🌟 What We've Built

AfriLearn is a **production-grade curriculum backend + LLM system prompt engine** purpose-built for African educational AI. Here is everything in the system:

### Core Infrastructure

| Component | File(s) | Description |
|---|---|---|
| **HTTP API Server** | `cmd/api/main.go` | Gin-based REST API with middleware, rate limiting, API key auth |
| **Database Layer** | `internal/database/` | PostgreSQL connection pool with `pgvector` extension support |
| **Schema & Migrations** | `cmd/migrate/main.go` | Auto-deploys all tables, GIN FTS indexes, PGVector HNSW column |
| **Curriculum Ingestion Engine** | `internal/ingestion/engine.go` | Reads 56 JSON files → upserts to PostgreSQL (idempotent) |
| **Data Seeder** | `cmd/seeder/main.go` | Runs ingestion for all datasets with `--validate-only` flag |
| **In-Memory Cache** | `internal/cache/` | Thread-safe LRU cache with per-prefix hit/miss metrics |
| **API Key Management** | `internal/handlers/keys.go` | Key generation, validation, and rate-limit metering |

### Curriculum API Endpoints

| Component | File(s) | Description |
|---|---|---|
| **Curriculum CRUD** | `internal/handlers/curriculum.go` | `GET /api/v1/curriculum/:board/:subject` — full topic tree |
| **LLM System Prompt Engine** | `internal/handlers/llm_prompt.go` | `GET /api/v1/curriculum/:board/:subject/llm-prompt` — domain-directed prompt |
| **RAG Embedding Chunker** | `internal/handlers/embeddings.go` | `GET /api/v1/curriculum/:board/:subject/embeddings` — 200–800 token chunks |
| **3-Layer Full-Text Search** | `internal/handlers/search.go` | `GET /api/v1/search` — ranked FTS across topics, subtopics, objectives |
| **Vector Search Endpoint** | `internal/handlers/vector.go` | `POST /api/v1/search/vector` — PGVector HNSW (Phase 2: ML embeddings) |
| **Cross-Board Matcher** | `internal/handlers/match.go` | `GET /api/v1/curriculum/match/:topic` — queries all 22 boards simultaneously |
| **Learning Pathway Engine** | `internal/handlers/pathway.go` | `GET /api/v1/curriculum/pathway/:topic` — BECE→WAEC→JAMB→Uni progression |
| **Prerequisite Graph** | `internal/handlers/prerequisites.go` | `GET /api/v1/curriculum/prerequisites/:board/:subject/:topic` |
| **Natural Language Query Brain** | `internal/handlers/query.go` | Intent parser with grade-level and topic inference fallbacks |

### Intelligence & Pedagogy Layer

| Component | File(s) | Description |
|---|---|---|
| **Domain Directive Engine** | `internal/handlers/llm_prompt.go` | Subject-specific rules: IRAC (Law), SOAP (Medicine), GIVEN/REQUIRED (Engineering), misconception flags (Physics/Chemistry) |
| **Bloom's Taxonomy Classifier** | `internal/handlers/llm_prompt.go` | Auto-classifies every topic across 6 cognitive levels |
| **Misconception Interception** | `internal/handlers/llm_prompt.go` | Explicit flags for zero-gravity myth, Lamarck confusion, condensation errors, Le Chatelier pressure traps |
| **Nigerian Context Anchoring** | `internal/handlers/llm_prompt.go` | Lagos analogies, NBS data hooks, Nigerian case law directives, WAEC/JAMB exam format rules |
| **Adaptive Learning Path** | `internal/handlers/llm_prompt.go` | Difficulty progression suggestions per curriculum level |

### Operations & DevOps

| Component | File(s) | Description |
|---|---|---|
| **Admin Dashboard** | `internal/handlers/admin.go` + `cmd/api/main.go` | `/admin` — real-time cache stats, purge, background re-ingestion |
| **Developer Portal** | `cmd/api/main.go` | `/` — styled landing with all endpoint links |
| **Interactive AI Playground** | `cmd/api/main.go` | `/playground` — live AI tutor chat testing (Groq/Gemini) |
| **Swagger / OpenAPI Docs** | `cmd/api/main.go` | `/docs` — interactive API explorer |
| **Load Test Utility** | `cmd/loadtest/main.go` | Concurrent benchmarking with p50/p90/p99 latency metrics |
| **Agentic Evaluation CLI** | `cmd/agentic_eval/main.go` | Parallel A/B test: Claude Sonnet 4.6 baseline vs AfriLearn activated |
| **Docker Support** | `Dockerfile`, `docker-compose.yml` | Containerized deployment (DB_URL injected via environment) |
| **Render Deployment** | `render.yaml` | One-click Render.com deployment config |

### Curriculum Datasets (56 files)

| Exam Board | Subjects Covered |
|---|---|
| **BECE** | Basic Science, Mathematics, English, Social Studies, Basic Technology |
| **WAEC** | Physics, Chemistry, Biology, Mathematics, English, Economics, Government, Literature |
| **JAMB** | Physics, Chemistry, Biology, Mathematics, English, Economics, Government, CRS |
| **UNILAG** | Law (LL.B.), Medicine & Surgery (MBBS), Computer Science, Economics, Accounting |
| **FUTO** | Mechanical Engineering, Computer Engineering, Electrical Engineering, Electronics |
| **NUC** | Computer Science (CCMAS), Electrical Engineering, Civil Engineering |
| **UNN / UNEC / EBSU / FUNAI** | Law, Computer Science, Business Administration |
| **YABATECH / IMT** | Computer Science ND/HND, Business Studies ND/HND |

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
Returns an LLM system prompt with Bloom's Taxonomy breakdown, domain-specific rules (Law, Medicine, Engineering, CS, etc.), difficulty progression, misconception interception flags, and token count.
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
Returns curriculum content pre-chunked (200–800 tokens per module) ready to embed with OpenAI, Gemini, or Ollama.
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
        "text_content": "CURRICULUM: West African Examinations Council — Physics\nLEVEL: SS1-SS3 | BOARD: WAEC\n...",
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

### 5. Vector Search Endpoint
```bash
curl -X POST -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  -H "Content-Type: application/json" \
  -d '{"query": "statutory interpretation rules in legal methods", "limit": 5}' \
  http://localhost:8080/api/v1/search/vector
```

### 6. Cross-Board Topic Matcher
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/match/thermodynamics
```

### 7. Learning Pathway Engine
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/pathway/electricity
```

---

## 🤖 LLM Benchmark: Claude Sonnet 4.6 vs AfriLearn System

AfriLearn includes a formal **agentic evaluation framework** that benchmarks the curriculum prompt engine against an unassisted **Claude Sonnet 4.6** baseline.

### Headline Result

| Metric | Claude Sonnet 4.6 (Raw) | AfriLearn System Activated | Gain |
|---|---|---|---|
| Section I (10 standard tests) | 54.6 / 100 | **97.6 / 100** | **+78.8%** |
| Section II (6 hard edge cases) | 45.0 / 100 | **95.8 / 100** | **+112.9%** |
| **Overall (16 scenarios)** | **50.9 / 100** | **96.9 / 100** | **+90.4%** |

> **Both Mode A and Mode B use the same Claude Sonnet 4.6 model.** The performance delta is entirely attributable to AfriLearn's curriculum data layer and domain directive engine — not model capability differences.

### What AfriLearn Adds Over Raw Claude

| Category | Raw Claude Sonnet 4.6 | AfriLearn Activated |
|---|---|---|
| **Exam format compliance** | Generic paragraphs | IRAC (Law), SOAP (Medicine), GIVEN/REQUIRED (Engineering) |
| **Misconception interception** | Answers question as asked | Leads with misconception confrontation box (❌ vs ✅) |
| **Nigerian context** | Global/generic examples | Lagos analogies, NBS data, Nigerian case law, WAEC mark schemes |
| **WAEC/JAMB exam tricks** | Not aware of exam patterns | Explicit JAMB Kc temperature trap warnings, WAEC no-lifting rule |
| **Graph pedagogy** | Text description only | ASCII graph sketches + fatal exam mistake warnings |

### Run the Evaluation

```bash
# Requires live DB connection
go run cmd/agentic_eval/main.go
```

Generates [`questionstest.md`](./questionstest.md) with full comparative responses and rubric audits.

**Detailed report**: [`BENCHMARK_REPORT.md`](./BENCHMARK_REPORT.md)

---

## ⚡ Performance Benchmarking & Load Testing

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

> [!CAUTION]
> Never commit `.env` or any file containing real database credentials to git. The `.gitignore` excludes `.env` by default. If credentials are exposed, rotate them immediately via your database provider dashboard.

---

## 🔒 Security Notes

- `.env` is git-ignored — credentials never enter version control
- `docker-compose.yml` uses `${DB_URL}` environment variable injection (no hardcoded secrets)
- API key validation middleware on all routes except `/health` and `/api/v1/keys/generate`
- In-memory rate limiting per API key (Redis-backed rate limiting planned for Phase 3)

---

## 📁 Project Structure

```
AfriLearn/
├── cmd/
│   ├── api/              # Main HTTP server entrypoint
│   ├── migrate/          # Schema deployment & migrations
│   ├── seeder/           # Curriculum ingestion runner
│   ├── loadtest/         # Concurrent load test utility
│   └── agentic_eval/     # A/B benchmark evaluation CLI
├── internal/
│   ├── cache/            # Thread-safe in-memory LRU cache
│   ├── database/         # PostgreSQL connection & query helpers
│   ├── handlers/         # All HTTP route handlers + LLM prompt engine
│   ├── ingestion/        # Curriculum JSON → DB engine
│   └── models/           # Shared data types & response models
├── data/
│   └── curricula/        # 56 JSON curriculum files (14 boards)
├── BENCHMARK_REPORT.md   # Claude Sonnet 4.6 vs AfriLearn formal report
├── questionstest.md      # 16-scenario A/B evaluation suite (1,147 lines)
├── ROADMAP.md            # 6-phase engineering roadmap
├── CONTRIBUTING.md       # Dataset contribution guide & JSON schema
├── Dockerfile            # Multi-stage Go container build
├── docker-compose.yml    # Local dev with environment variable injection
└── render.yaml           # One-click Render.com deployment
```

---

## 🗺️ Roadmap Summary

| Phase | Focus | Status |
|---|---|---|
| **Phase 1** | Core API + 56 curriculum datasets + FTS search | ✅ Complete |
| **Phase 2** | ML vector embeddings (OpenAI/Gemini → PGVector HNSW) + automated scoring LLM-as-judge | 🔄 In Progress |
| **Phase 3** | Student progress tracking + spaced repetition + Redis cache backend | 📅 Planned |
| **Phase 4** | Adaptive learning pathways + prerequisite graph + cross-board recommendations | 📅 Planned |
| **Phase 5** | hkai.site frontend integration + real-time data hooks (NBS, JAMB updates) | 📅 Planned |
| **Phase 6** | Multi-language support (Yoruba, Igbo, Hausa interface layer) | 📅 Planned |

Full details: [`ROADMAP.md`](./ROADMAP.md)

---

## ⚖️ License

The AfriLearn Curriculum API and all raw curriculum datasets in `data/curricula/` are distributed under the **MIT License**.

---

*Built for African students, by builders who care. Open source. Open data. Open future.*
