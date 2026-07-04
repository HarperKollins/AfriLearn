# AfriLearn Curriculum API

> **The foundational digital backbone & curriculum data layer for African educational technology, school portals, edtech apps, and AI Tutors.**  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB UTME · NUC CCMAS Degrees · NBTE Polytechnics · Federal & State Universities · Developer Portal (`/portal`) · Admin Dashboard (`/admin`) · In-Memory Caching · AI Tutor Prompt Engine · Interactive Swagger Docs (`/docs`) · Docker & Render Cloud Ready**

---

## 🌟 Key Features & Architecture

- **56 Standardized Curriculum Datasets**: Complete coverage across 14 exam boards and tertiary institutions including BECE, WAEC, JAMB, NECO, NUC, NBTE, YABATECH, IMT Enugu, UNILAG, UNN, UNEC, EBSU, FUNAI, and FUTO.
- **⚡ In-Memory Hot-Path Caching Layer**: Thread-safe `sync.RWMutex` cache for high-throughput responses on curriculum detail and AI Tutor prompt generation endpoints (`/curriculum` and `/llm-prompt`).
- **🛠️ Internal Admin Operations Dashboard (`/admin`)**: Interactive web portal to inspect live cache hit ratios, run dataset quality validation checks, purge memory cache, and trigger background data re-ingestion.
- **🛡️ Dataset Schema & Quality Validator**: Pre-ingestion validation pipeline (`cmd/seeder/main.go --validate-only`) verifying required fields, difficulty tags, and Bloom's verb structures before database writes.
- **🤖 AI Tutor System Prompt Engine**: Formats curriculum trees into structured system prompts with **Bloom's Taxonomy Analytics**, **Pedagogical Directives** (tailored to primary, secondary, and university levels), and **Adaptive Learning Paths**.
- **🔍 Intelligence Layer**: Cross-board curriculum alignment matcher, learning pathway generator (BECE → WAEC → JAMB → NUC), prerequisite graph engine, and natural language Query Brain.

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
AfriLearn uses PostgreSQL (compatible with Neon.tech, AWS RDS, Supabase, or local Postgres). Open `.env` and set your `DB_URL`:
```ini
DB_URL="postgresql://user:password@localhost:5432/afrilearn?sslmode=disable"
```

### 3. Run Pre-Ingestion Data Validation (Optional Quality Check)
```bash
go run cmd/seeder/main.go --validate-only
```

### 4. Deploy Database Schema & Ingest Curricula
```bash
# Deploy database schema, tables, GIN indexes, and seed reference data
go run cmd/migrate/main.go

# Ingest all 56 JSON curriculum datasets from data/curricula/
go run cmd/seeder/main.go
```

### 5. Start the API Server
```bash
go run cmd/api/main.go
```

- **Developer Self-Service Portal:** [http://localhost:8080/portal](http://localhost:8080/portal)
- **Admin & Operations Dashboard:** [http://localhost:8080/admin](http://localhost:8080/admin)
- **Interactive Swagger API Docs:** [http://localhost:8080/docs](http://localhost:8080/docs)
- **Health Check:** [http://localhost:8080/health](http://localhost:8080/health)

---

## 🛠️ Admin & Operations Dashboard (`/admin`)

Visit `http://localhost:8080/admin` to access the internal management portal:
- **Live Cache Metrics**: Monitor total cached items, hits, misses, and hit ratio.
- **⚡ Re-Ingest All Curricula**: Trigger background scanning and upserting of all 56 curriculum files directly from the UI.
- **🔍 Dataset Validation**: Run real-time schema & objective quality checks across all JSON files.
- **🧹 Purge Memory Cache**: Instantly flush cached API responses.

---

## 🔑 Developer Portal & Self-Service API Keys (`/portal`)

Open `http://localhost:8080/portal` to manage API keys:
- Generate instant self-service API keys for **Free Tier** (1,000 req/min) or **Pro Partner Tier** (50,000 req/min).
- Built-in live API testing console with copyable curl examples.
- Non-blocking batched memory metering worker that flushes usage stats to PostgreSQL periodically without database locks.

---

## 🤖 AI Tutor Prompt Engine & Bloom's Analytics (`/llm-prompt`)

AfriLearn formats curriculum trees into system prompts ready for LLM context windows (GPT-4, Claude, Gemini, LLaMA):
```bash
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" http://localhost:8080/api/v1/curriculum/unilag/law/llm-prompt
```

### Response Features:
- **`system_prompt`**: Formatted prompt defining tutor persona and academic depth.
- **`blooms_taxonomy_breakdown`**: Cognitive analytics (`remember`, `understand`, `apply`, `analyze`, `evaluate`, `create`).
- **`pedagogical_directives`**: Level-specific teaching guidelines (e.g. legal case law citations for UNILAG Law, clinical case presentations for UNILAG MBBS, practice drills for WAEC).
- **`adaptive_learning_path`**: Ordered step-by-step module mastery sequence with difficulty tags (`[EASY]`, `[MEDIUM]`, `[HARD]`).

---

## 🔍 Intelligence Layer Endpoints

1. **Cross-Board Curriculum Matcher (`/api/v1/curriculum/match/:topic`)**  
   Queries all boards simultaneously using PostgreSQL Full-Text Search (GIN FTS index), returning how a topic is taught across junior secondary, senior secondary, and tertiary degree levels.
2. **Learning Pathway Engine (`/api/v1/curriculum/pathway?subject=...&from=bece&to=jamb`)**  
   Constructs a step-by-step learning progression journey mapping a student's path from JSS1 to UTME or University.
3. **Prerequisite Engine (`/api/v1/curriculum/prerequisites/:board/:subject/:topic`)**  
   Queries the dependency graph to return foundational topics a student must master before attempting the target topic.
4. **Natural Language Query Brain (`/api/v1/query`)**  
   Accepts natural language user queries (e.g. *"How do I study law at UNILAG?"*), parses intent, handles multi-turn clarification loops, and caches responses.

---

## ✍️ Data Decoupling & Adding New Curricula

Curricula are completely decoupled from Go code. Adding a subject or institution requires **no code changes or recompilation**.

### Directory Structure (`data/curricula/`):
- `data/curricula/bece/` — Junior Secondary (JSS1-3)
- `data/curricula/waec/` — Senior Secondary (SS1-3)
- `data/curricula/jamb/` — UTME Entrance
- `data/curricula/nuc/` — National Universities Commission (CCMAS)
- `data/curricula/unilag/`, `unn/`, `unec/`, `ebsu/`, `funai/`, `futo/` — University Degree Programs
- `data/curricula/yabatech/`, `imt/` — Polytechnic ND/HND Standards

Save your file as `data/curricula/[board]/[subject].json` and run `go run cmd/seeder/main.go`.

---

## 🐳 Containerization & Cloud Deployment

### Run with Docker Compose
```bash
docker-compose up -d --build
```

### Deploy to Render / Railway / AWS
The repository includes a production-ready `Dockerfile` and `render.yaml`:
1. Connect your GitHub repository to Render.
2. Render auto-detects the blueprint and provisions the Go web service.
3. Set `DB_URL` in your deployment environment variables.
