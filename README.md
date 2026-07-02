# AfriLearn Curriculum API

> The foundational data layer for African educational technology & AI Tutors.  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB · NUC CCMAS (100L-500L) · NBTE Polytechnics (ND/HND) · Developer API Keys · AI Tutor LLM Prompts · Interactive Swagger UI (`/docs`) · Docker & Cloud Ready** — structured as APIs.

---

## What This Is

An infrastructure API that exposes official Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, universities, polytechnics, schools, and educational platforms.

```bash
# 📖 Interactive Swagger UI Playground
http://localhost:8080/docs

# Authenticated Request with X-API-Key
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" http://localhost:8080/api/v1/curriculum/waec/physics

# 🤖 AI Tutor LLM System Prompt & Context Window Endpoint
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" http://localhost:8080/api/v1/curriculum/waec/physics/llm-prompt
```

---

## Quick Start

### Local Go Setup
```bash
git clone https://github.com/HarperKollins/AfriLearn.git
cd eduscrape
go mod tidy
cp .env.example .env
go run cmd/migrate/main.go
go run cmd/seeder/main.go
go run cmd/api/main.go
# Open http://localhost:8080/docs in your browser
```

---

## 🐳 Docker & Cloud Production Deployment

### 1. Run with Docker Compose
```bash
docker-compose up -d --build
```

### 2. Build & Run Docker Image Manually
```bash
docker build -t afrilearn-api .
docker run -p 8080:8080 -e DB_URL="<your_neon_db_url>" afrilearn-api
```

### 3. 1-Click Deploy on Render.com
1. Connect your GitHub repository (`HarperKollins/AfriLearn`) to [Render.com](https://render.com).
2. Render automatically detects `render.yaml` and deploys the web service with free SSL!
3. Add your `DB_URL` environment variable under Render dashboard settings.

---

## Developer API Keys & Rate Tiers

Pass your API key in the `X-API-Key` HTTP header or as an `api_key` query parameter.

| Tier | API Key | Rate Limit | Description |
|------|---------|------------|-------------|
| **Public** | *(None)* | 60 req/min | Public access without header |
| **Free** | `afr_live_demo_9f8e2b7a` | 1,000 req/min | Free Developer Key |
| **Pro** | `afr_live_pro_8372bf91` | 50,000 req/min | Commercial EdTech Partner Key |

---

## Live API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | API health check & DB status |
| GET | `/docs` | **📖 Interactive Swagger UI Playground** (OpenAPI 3.0) |
| GET | `/docs/openapi.json` | OpenAPI 3.0 JSON specification |
| GET | `/api/v1/` | API info and endpoint index |
| GET | `/api/v1/subjects` | List all 46 subjects, university degrees & polytechnic diploma programs |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all 22 exam boards, polytechnics & universities |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum tree with topics, subtopics & objectives |
| GET | `/api/v1/curriculum/:board/:subject/llm-prompt` | **🤖 AI Tutor Endpoint**: Formatted System Prompt, Context Window, and Module Blocks for GPT-4, Claude, Gemini, and LLaMA |
| GET | `/api/v1/search?q=:query` | Search topics across all exam boards, university degrees, and polytechnic diplomas |

---

## Technical Architecture

- **Go (Golang)** — High-performance, low-latency compiled runtime
- **Gin** — Fast HTTP framework with CORS & middleware
- **PostgreSQL (Neon)** — Relational storage with JSON array batching optimization (`pq.Array`)
- **API Key Auth** — In-memory cached key authentication & background usage metering (`internal/middleware/auth.go`)
- **OpenAPI 3.0 Docs** — Embedded Swagger UI documentation playground (`internal/handlers/docs.go`)
- **Docker Multi-Stage** — Minimal Alpine production container (`Dockerfile`, `docker-compose.yml`, `render.yaml`)
- **Scraper Engine** — Modular `Scraper` interface pattern (`internal/scraper`)
