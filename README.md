# AfriLearn Curriculum API

> **The foundational digital backbone & curriculum data layer for African educational technology, school portals, edtech apps, and AI Tutors.**
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB UTME · NUC CCMAS Degrees · NBTE Polytechnics · Open Data Community · AI Tutor LLM Prompts · Interactive Swagger Docs**

---

## 🌟 What This Is

AfriLearn is a **structured curriculum database API** covering **56 official Nigerian educational datasets** — from junior secondary school all the way through university degrees. It is designed to be the data layer that powers AI tutors, school portals, and edtech applications.

### Current Capabilities (Honest)

| Feature | Status | Notes |
|---|---|---|
| 56 structured curriculum datasets | ✅ Production | BECE, WAEC, JAMB, NUC, NBTE + 14 universities |
| Full-text search (topics + subtopics + objectives) | ✅ Production | PostgreSQL GIN, 3-layer ranked search |
| LLM system prompt engine | ✅ Production | Bloom's Taxonomy, subject-specific rules (Law, Medicine, Engineering...) |
| Cross-board curriculum matcher | ✅ Production | `GET /api/v1/curriculum/match/:topic` |
| Learning pathway engine | ✅ Production | BECE → WAEC → JAMB → University ordering |
| Prerequisite graph | ✅ Production | Explicit + derived prerequisites |
| Natural language query brain | ✅ Production | Intent parser + clarification loop |
| Admin operations dashboard | ✅ Production | `/admin` — cache stats, purge, re-ingest |
| Interactive AI Playground | ✅ Production | `/playground` — test all endpoints + AI chat |
| Swagger/OpenAPI docs | ✅ Production | `/docs` |
| In-memory caching (hot paths) | ✅ Production | Thread-safe, per-endpoint key |
| PGVector HNSW index | ✅ Schema-ready | Provisioned, ready for real embeddings |
| RAG text chunking (`/embeddings`) | ✅ Production | Correctly formatted chunks for any embedding model |
| **Real semantic vector search** | 🔄 Phase 2 | Needs embedding model integration — see ROADMAP.md |

> [!NOTE]
> **On the vector search**: The `/api/v1/search/vector` endpoint uses **PostgreSQL full-text search** — not ML vector embeddings. The pgvector HNSW index is provisioned and schema-ready. The `/embeddings` endpoint provides correctly formatted text chunks for you to embed with OpenAI/Gemini/Ollama. See [ROADMAP.md](./ROADMAP.md) for the integration guide.

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

### 4. Run Tests
```bash
go test -v ./...
```

### 5. Start the API Server
```bash
go run cmd/api/main.go
```

| URL | What it is |
|---|---|
| `http://localhost:8080/` | Developer portal |
| `http://localhost:8080/admin` | Admin operations dashboard |
| `http://localhost:8080/playground` | Interactive AI Playground |
| `http://localhost:8080/docs` | Swagger/OpenAPI interactive docs |
| `http://localhost:8080/health` | Health check |

---

## 🔑 Authentication

All API endpoints (except health and key generation) require an `X-API-Key` header.

**Demo key** (rate-limited, for testing): `afr_live_demo_9f8e2b7a`

**Generate a real key** (unlimited, self-service):
```bash
curl -X POST http://localhost:8080/api/v1/keys/generate \
  -H "Content-Type: application/json" \
  -d '{"name": "My App", "email": "dev@example.com"}'
```

---

## 📚 Core API Endpoints

### Get Full Curriculum
Returns a board's complete curriculum tree: topics → subtopics → learning objectives.
```bash
# WAEC Mathematics
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/waec/mathematics

# UNILAG Law (LL.B.)
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/unilag/law

# FUTO Engineering
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/futo/mechanical-engineering
```

### Get AI Tutor System Prompt
Returns a complete LLM system prompt with Bloom's Taxonomy breakdown, subject-specific pedagogical rules, difficulty progression, and token count estimate.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/unilag/law/llm-prompt
```

**Key fields in response**:
```json
{
  "system_prompt": "You are an expert AI Tutor for UNILAG LL.B. Law...",
  "subject_specific_rules": [
    "Apply the IRAC method for problem questions: Issue → Rule → Application → Conclusion.",
    "Always cite the full case name, court, and year (e.g., Donoghue v Stevenson [1932] AC 562 (HL)).",
    "Reference the Nigerian Constitution 1999 and CAMA 2020 by section number."
  ],
  "blooms_taxonomy_breakdown": {"remember": 12, "understand": 34, "apply": 28, "analyze": 15},
  "estimated_token_count": 4200,
  "suggested_chunking_note": "~4200 tokens — fits in most 8K models. Use full_context_window directly.",
  "difficulty_progression": ["BEGINNER", "INTERMEDIATE", "INTERMEDIATE", "ADVANCED"]
}
```

### Search Across All Curricula
Searches topics, subtopics, AND learning objectives simultaneously with ranked results.
```bash
# Basic search
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/search?q=photosynthesis"

# Filtered search (by board and subject) with pagination
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/search?q=newton+laws&board=waec&subject=physics&limit=10&offset=0"
```

**Response shows which layer matched**:
```json
{
  "data": [
    {
      "topic_name": "Force, Motion and Energy",
      "board_name": "WAEC",
      "subject_name": "Physics",
      "matched_in": "subtopic",
      "snippet": "Newton's Laws of Motion",
      "relevance_score": 0.0759
    }
  ]
}
```

### Get RAG Embedding Chunks
Returns curriculum content pre-chunked (one chunk per topic module) for embedding into any vector database. Pass `text_content` through your embedding model of choice.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/waec/physics/embeddings
```

**Response includes integration guide**:
```json
{
  "chunking_strategy": "one-chunk-per-topic-module (200–800 tokens per chunk)",
  "embedding_note": "Pass text_content to your embedding model (OpenAI, Gemini, Ollama)...",
  "integration_guide": {
    "openai": "client.embeddings.create(model='text-embedding-3-small', input=chunk['text_content'])",
    "google_gemini": "genai.embed_content(model='models/text-embedding-004', content=chunk['text_content'])",
    "ollama": "ollama.embeddings(model='nomic-embed-text', prompt=chunk['text_content'])"
  }
}
```

---

## 🧠 Intelligence Layer

### Cross-Board Curriculum Matcher
Find how a topic is taught across ALL Nigerian curriculum levels simultaneously.
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/match/quadratic-equations
```
Returns coverage across BECE → WAEC → JAMB → NUC → all universities, with a unified LLM prompt spanning all levels.

### Learning Pathway Engine
Get an ordered learning journey through curriculum levels for a subject.
```bash
# Full journey for Mathematics (BECE through to University)
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/curriculum/pathway?subject=mathematics"

# Scoped journey (BECE to JAMB only)
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  "http://localhost:8080/api/v1/curriculum/pathway?subject=mathematics&from=bece&to=jamb"
```

### Natural Language Query Brain
Ask curriculum questions in plain English — the system parses intent and routes to the right data.
```bash
curl -X POST -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  -H "Content-Type: application/json" \
  -d '{"query": "show me WAEC physics topics on electromagnetism"}' \
  http://localhost:8080/api/v1/query
```

---

## 🗃️ Supported Boards & Institutions (22 total)

| Slug | Institution | Level |
|---|---|---|
| `bece` | NERDC / Junior WAEC | JSS1–3 |
| `waec` | West African Examinations Council | SS1–3 |
| `neco` | National Examinations Council | SS1–3 |
| `jamb` | Joint Admissions & Matriculation Board | UTME |
| `nuc` | National Universities Commission (national standard) | Degree |
| `nbte` | National Board for Technical Education | ND/HND |
| `unilag` | University of Lagos | Degree |
| `unn` | University of Nigeria, Nsukka | Degree |
| `unec` | University of Nigeria, Enugu Campus | Degree |
| `ebsu` | Ebonyi State University | Degree |
| `funai` | Federal University, Ndufu-Alike Ikwo | Degree |
| `futo` | Federal University of Technology, Owerri | Degree |
| `yabatech` | Yaba College of Technology | ND/HND |
| `imt` | Institute of Management & Technology, Enugu | ND/HND |

---

## 🛠️ Admin Operations

The admin dashboard at `/admin` provides:
- **Real-time cache statistics** — hit rate, entry count, memory usage
- **Cache purge** — clear all cached responses
- **Background re-ingestion** — re-seed the database from disk without stopping the server
- **Dataset validation** — run schema checks on all 56 JSON files

API endpoints (no auth required — for internal use):
```bash
GET  /api/v1/admin/cache/stats     # Cache hit/miss stats
POST /api/v1/admin/cache/purge     # Clear cache
POST /api/v1/admin/reingest        # Trigger background re-ingestion
GET  /api/v1/admin/validate        # Validate all curriculum datasets
```

---

## 🌐 Open Data & Community

All curriculum data in `data/curricula/` is **100% open-source under the MIT License**. We welcome contributions from teachers, professors, and developers across Africa.

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for data submission guidelines, JSON schema, and PR workflows.

**Priority contributions needed**: Pharmacy (B.Pharm), Architecture (B.Arch), Agricultural Science, NECO-specific datasets, 400-level topics for existing university programs.

---

## 🗺️ Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for the full technical roadmap including:
- **Phase 2**: Real semantic vector search (OpenAI/Gemini/Ollama embedding integration)
- **Phase 3**: Student progress tracking
- **Phase 4**: Adaptive learning pathways
- **Phase 5**: Multi-modal (video, audio, OCR of Nigerian textbooks)
- **Phase 6**: National infrastructure layer for Nigerian edtech

---

## ⚙️ Environment Variables

```ini
DB_URL=postgresql://user:password@host:5432/afrilearn?sslmode=require
PORT=8080
APP_ENV=development   # or: production
```

See [`.env.example`](./.env.example) for the full reference.

---

## 📦 Docker

```bash
docker-compose up
```

The service runs on port `8080`. Ensure `DB_URL` is set in `docker-compose.yml`.
