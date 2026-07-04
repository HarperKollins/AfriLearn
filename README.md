# AfriLearn Curriculum API

> **The foundational digital backbone & vector data primitive for African educational technology, school portals, edtech apps, and AI Tutors.**  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB UTME · NUC CCMAS Degrees · NBTE Polytechnics · PGVector HNSW Cosine Search · Open Data Community (`CONTRIBUTING.md`) · AI Tutor LLM Prompts · Interactive Swagger Docs (`/docs`) · Cloud Ready**

---

## 🌟 Key Architecture & Capabilities

- **56 Standardized Curriculum Datasets**: Complete coverage across 14 exam boards and tertiary institutions including BECE, WAEC, JAMB, NECO, NUC, NBTE, YABATECH, IMT Enugu, UNILAG, UNN, UNEC, EBSU, FUNAI, and FUTO.
- **⚡ Native PostgreSQL PGVector HNSW Engine (`POST /api/v1/search/vector`)**:
  - **Vector Dimension**: 1536-dimensional L2-normalized float embeddings.
  - **Index Specs**: HNSW (Hierarchical Navigable Small World) index with `m=16, ef_construction=64` for sub-5ms cosine similarity search.
  - **Search Query Mode**: `ORDER BY embedding <=> $1::vector` returning cosine distance matching scores.
- **🧠 Granular Subtopic & Objective RAG Chunking (`GET /api/v1/curriculum/:board/:subject/embeddings`)**: Pre-chunks curriculum modules into 200-400 token semantic blocks tailored for LLM context injection into Pinecone, PGVector, Qdrant, or ChromaDB.
- **⚡ Thread-Safe In-Memory Caching Layer**: Multi-level cache with sub-millisecond response times for hot endpoints (`/curriculum` and `/llm-prompt`).
- **🛠️ Operations Dashboard (`/admin`)**: Web portal to monitor real-time cache hit ratios, run dataset quality validation checks, purge memory cache, and trigger background data re-ingestion.
- **🛡️ Dataset Schema Validator**: CLI tool (`cmd/seeder/main.go --validate-only`) verifying required fields, difficulty tags, and Bloom's verb structures.
- **🤖 AI Tutor System Prompt Engine**: Formats curriculum trees into structured system prompts with **Bloom's Taxonomy Analytics**, **Pedagogical Directives**, and **Adaptive Learning Paths**.

---

## 🚀 Quick Start

### 1. Clone & Set Up
```bash
# Clone the repository
git clone https://github.com/HarperKollins/AfriLearn.git
cd AfriLearn

# Install dependencies
go mod tidy

# Create environment configuration
cp .env.example .env
```

### 2. Configure Database
AfriLearn uses PostgreSQL with native `pgvector` extension (Neon.tech, AWS RDS, Supabase, or local Postgres). Set your `DB_URL` in `.env`:
```ini
DB_URL="postgresql://user:password@localhost:5432/afrilearn?sslmode=disable"
```

### 3. Deploy Schema & Ingest Curricula
```bash
# Deploys schema, tables, GIN full-text indexes, and PGVector HNSW index
go run cmd/migrate/main.go

# Ingests all 56 curriculum files and computes 1536-dim vector embeddings
go run cmd/seeder/main.go
```

### 4. Run Automated Test Suite
```bash
go test -v ./...
```

### 5. Start the API Server
```bash
go run cmd/api/main.go
```

- **Developer Portal:** [http://localhost:8080/portal](http://localhost:8080/portal)
- **Admin Dashboard:** [http://localhost:8080/admin](http://localhost:8080/admin)
- **Interactive Swagger Docs:** [http://localhost:8080/docs](http://localhost:8080/docs)

---

## ⚡ Native PGVector Cosine Similarity Search (`POST /api/v1/search/vector`)

Send natural language text queries to query topics by 1536-dim vector cosine similarity:

```bash
curl -X POST -H "X-API-Key: afr_live_demo_9f8e2b7a" \
     -H "Content-Type: application/json" \
     -d '{"query": "statutory interpretation rules in legal methods", "limit": 5}' \
     http://localhost:8080/api/v1/search/vector
```

### Response Payload:
```json
{
  "success": true,
  "data": {
    "query": "statutory interpretation rules in legal methods",
    "search_mode": "pgvector-cosine-similarity (HNSW Index)",
    "vector_dimension": 1536,
    "total_matches": 5,
    "results": [
      {
        "board_name": "UNILAG",
        "subject_name": "LL.B. Bachelor of Laws",
        "topic_name": "Year 1 (100 Level): Legal Foundations & Methodologies",
        "similarity_score": 1.0000000000000004
      },
      {
        "board_name": "NUC",
        "subject_name": "Law (LL.B.)",
        "topic_name": "LAW 101: Legal Methods & Nigerian Legal System",
        "similarity_score": 0.9329712613149811
      }
    ]
  }
}
```

---

## 🌐 Open Data & Community Contributions

Raw curriculum data in `data/curricula/` is 100% open-source under the MIT License. We welcome contributions from teachers, professors, and developers!

See [`CONTRIBUTING.md`](file:///c:/Users/Harper/Desktop/eduscrape/CONTRIBUTING.md) for data submission guidelines, JSON schema rules, and PR workflows.
