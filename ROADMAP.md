# AfriLearn — Technical Roadmap

> **Current status**: Production MVP — 56 curriculum datasets, FTS search, LLM prompt engine, admin dashboard, intelligence layer (cross-board matcher, pathway engine, prerequisite graph, query brain).

---

## Phase 1: Data Foundation ✅ Complete

*Goal: Build the canonical, structured curriculum database for Nigeria.*

| Milestone | Status |
|---|---|
| 56 curriculum datasets (BECE, WAEC, JAMB, NUC, NBTE, 14 institutions) | ✅ Done |
| PostgreSQL schema with GIN full-text search indexes | ✅ Done |
| Idempotent ingestion engine (zero duplicate objectives across runs) | ✅ Done |
| Deterministic subtopic ordering (order_index, not insertion order) | ✅ Done |
| Dataset schema validator (`--validate-only` CLI flag) | ✅ Done |
| REST API with auth, rate limiting, CORS | ✅ Done |
| PGVector HNSW index provisioned (schema-ready for embeddings) | ✅ Done |
| In-memory caching (hot path: curriculum + LLM prompt) | ✅ Done |
| Admin operations dashboard (`/admin`) | ✅ Done |
| Interactive AI Playground (`/playground`) | ✅ Done |
| Swagger/OpenAPI docs (`/docs`) | ✅ Done |
| Intelligence layer: cross-board matcher, pathway engine, prerequisite graph | ✅ Done |
| Natural language query brain (`POST /api/v1/query`) | ✅ Done |

**Data gaps to fill (Phase 1 ongoing)**:
- [ ] More tertiary programs: Pharmacy, Architecture, Agricultural Science, Public Administration
- [ ] 400-level topics for existing university datasets (most currently cover 100–300 level)
- [ ] NECO-specific datasets (currently using WAEC where they overlap)
- [ ] Federal Polytechnic datasets (MAPOLY, Fed Poly Ilaro, Nekede)

---

## Phase 2: Real Semantic Intelligence 🔄 Next

*Goal: Replace the FTS search with true semantic vector search using real ML embeddings.*

**The Problem**: The current `/api/v1/search/vector` and `/api/v1/curriculum/:board/:subject/embeddings` use deterministic hash-based placeholder vectors. They are structurally sound (pgvector HNSW index is provisioned, chunk formatting is correct) but are semantically meaningless — two topics on "quadratic equations" will have a random cosine distance, not a meaningful semantic similarity.

**The Solution**: Integrate a real embedding model. The `/embeddings` endpoint already produces correctly chunked text — you just need to pass `text_content` through an embedding model and upsert into `topics.embedding`.

### Implementation Plan

**Option A: OpenAI (Highest quality, $0.02/1M tokens)**
```python
from openai import OpenAI
client = OpenAI()

# Get chunks from API
chunks = requests.get('/api/v1/curriculum/waec/mathematics/embeddings').json()['data']['chunks']

for chunk in chunks:
    embedding = client.embeddings.create(
        model='text-embedding-3-small',
        input=chunk['text_content']
    ).data[0].embedding
    
    # Upsert into pgvector
    db.execute(
        "UPDATE topics SET embedding = $1::vector WHERE id IN "
        "(SELECT id FROM topics WHERE slug = $2 LIMIT 1)",
        [embedding, chunk['chunk_id']]
    )
```

**Option B: Google Gemini (Free tier available)**
```python
import google.generativeai as genai
genai.configure(api_key='YOUR_KEY')

result = genai.embed_content(
    model='models/text-embedding-004',
    content=chunk['text_content']
)
embedding = result['embedding']  # 768-dim, update pgvector column accordingly
```

**Option C: Ollama / Local (Free, private, self-hosted)**
```bash
ollama pull nomic-embed-text  # 768-dim
ollama pull mxbai-embed-large  # 1024-dim
```
```python
import ollama
embedding = ollama.embeddings(model='nomic-embed-text', prompt=chunk['text_content'])['embedding']
```

### What needs to change in the codebase
1. **`cmd/seeder/main.go`** — add `--embed` flag that reads `EMBEDDING_PROVIDER` (openai/gemini/ollama) and `EMBEDDING_API_KEY` from `.env`, calls the provider, and upserts into `topics.embedding`
2. **`internal/handlers/embeddings.go`** — remove `EmbeddingModel: "placeholder"` and replace with dynamic field based on what's stored
3. **`internal/handlers/embeddings.go`** — `HandleVectorSearch` already has pgvector search in Phase 1 — it will automatically work once real embeddings are stored
4. **Migration** — optionally resize the `embedding` column to match the chosen model dimension (currently 1536 for OpenAI, may need 768 for Gemini/Ollama)

**Estimated effort**: 1–2 days  
**Cost estimate**: ~$0.50 for OpenAI to embed all 56 curricula × ~50 topics each = ~2,800 chunks × ~300 tokens/chunk = ~840K tokens = ~$0.017 total  

---

## Phase 3: Student Progress Tracking 📋 Planned

*Goal: Connect to hkai.site frontend — track what students have studied, mastered, and struggled with.*

### Schema additions needed
```sql
-- Student progress table
CREATE TABLE student_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,         -- FK to users table on hkai.site
    topic_id UUID REFERENCES topics(id),
    subtopic_id UUID REFERENCES subtopics(id),
    status VARCHAR(20) DEFAULT 'not_started',  -- not_started, in_progress, completed, struggling
    mastery_score DECIMAL(5,2),        -- 0-100 mastery percentage
    attempts INT DEFAULT 0,
    last_studied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, subtopic_id)
);

-- Student sessions
CREATE TABLE tutor_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    board_slug VARCHAR(50),
    subject_slug VARCHAR(50),
    messages JSONB,                    -- full conversation history
    session_start TIMESTAMPTZ DEFAULT NOW(),
    session_end TIMESTAMPTZ
);
```

### New API endpoints
```
POST /api/v1/student/progress          -- Update topic/subtopic mastery
GET  /api/v1/student/:id/progress      -- Get student's full progress map
GET  /api/v1/student/:id/weak-areas    -- Topics with mastery < 60%
GET  /api/v1/student/:id/next-topic    -- Adaptive: suggest what to study next
```

---

## Phase 4: Adaptive Learning Pathways 🧠 Planned

*Goal: ML-driven personalized learning path recommendations based on student performance.*

**How it works**:
1. When a student completes a subtopic with mastery < 70%, suggest prerequisite topics
2. When mastery ≥ 90%, unlock next topic in the pathway automatically
3. Use the existing prerequisite graph + cross-board pathway engine to route students
4. A/B test different pathway orderings to find optimal sequences

### Adaptive recommendation algorithm
```
Given: student S has mastered topics T1, T2 and struggled with T3
Find: next_topic = argmax_{t ∉ mastered} [
    0.4 × cos_similarity(t.embedding, T3.embedding)  -- topically related to struggle
  + 0.4 × (prerequisite_coverage(t) / prerequisite_count(t))  -- prerequisites met
  + 0.2 × (1 / topic.difficulty_score)  -- prefer easier topics when struggling
]
```

---

## Phase 5: Multi-Modal Expansion 🎥 Vision

*Goal: Go beyond text — support video lectures, audio notes, and scanned Nigerian textbook OCR.*

| Feature | Description |
|---|---|
| **Video curriculum** | Link topics to YouTube/Wasilisha video IDs |
| **Audio summaries** | TTS-generated topic summaries for offline/low-bandwidth use |
| **Nigerian textbook OCR** | Scan NERDC-approved textbooks → extract topics → ingest |
| **Past question bank** | Structure WAEC/JAMB past questions → link to topics |
| **Teacher mode** | Let teachers annotate curricula with local examples |

---

## Phase 6: National Infrastructure Play 🌍 Long-term Vision

*The AfriLearn curriculum API isn't just for hkai.site — it's a public infrastructure layer for African edtech.*

**Goal**: Be to Nigerian education what Stripe is to payments — the foundational API that every edtech startup builds on.

```
Current use cases being explored:
├── School management systems (curriculum compliance checking)
├── Content authoring tools (structured lesson plan generation)
├── Assessment platforms (map questions to official syllabi)
├── Parent apps (track children's curriculum coverage)
├── Teacher training (generate lesson plans aligned to NUC/NERDC)
└── Government portals (UBEC, SAME, state education ministries)
```

**What this requires**:
- SLA commitments (99.9% uptime)
- Paid API tier with usage metering (Stripe integration)
- Versioned APIs (v1, v2) with deprecation notices
- Redis/DB-backed cache (not in-memory, for multi-instance deployment)
- Rate limiting by organization, not just IP
- Audit logs for government/enterprise clients

---

## Known Limitations (Honest)

| Limitation | Impact | Mitigation |
|---|---|---|
| In-memory cache wipes on restart | Hot data cold again after deploy | Phase 2: Redis or DB-backed cache |
| Hash-based "vector" embeddings | Vector search = FTS in disguise | Phase 2: Real embedding model |
| 56 datasets, depth varies | Some tertiary programs surface-level | Phase 1 ongoing: community contributions |
| No freshness checks vs. official sources | Curricula may drift from official | Quarterly diff reviews against official syllabi |
| Single-instance rate limiter | Doesn't work with horizontal scaling | Redis-backed rate limiter for scale |

---

## Contributing to Phase 1 (Right Now)

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for:
- JSON schema for curriculum datasets
- How to research and submit a new dataset
- Data quality validation checklist
- PR review process

Priority datasets needed:
1. **Pharmacy (B.Pharm)** — UNILAG, UNIBEN, OAU
2. **Architecture (B.Arch)** — NUC standard
3. **Agricultural Science (B.Agric)** — OAU, FUNAAB, ABU
4. **Mass Communication (100–400 level)** — complete depth for existing boards
5. **NECO-specific** — differentiate from WAEC where syllabus diverges
