# AfriLearn Curriculum API

> The foundational data layer for African educational technology & AI Tutors.  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB · NUC CCMAS (100L-500L) · NBTE Polytechnics (ND/HND) · AI Tutor LLM Prompt Generation** — structured as APIs.

---

## What This Is

An infrastructure API that exposes official Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, universities, polytechnics, schools, and educational platforms.

```
# 🤖 AI Tutor LLM System Prompt & Context Window Endpoints
GET /api/v1/curriculum/waec/physics/llm-prompt
GET /api/v1/curriculum/jamb/mathematics/llm-prompt
GET /api/v1/curriculum/nuc/computer-science/llm-prompt
GET /api/v1/curriculum/yabatech/computer-engineering-tech/llm-prompt

# Full Curriculum Tree Endpoints
GET /api/v1/curriculum/yabatech/computer-engineering-tech
GET /api/v1/curriculum/unilag/computer-science
GET /api/v1/curriculum/nuc/computer-science
GET /api/v1/curriculum/waec/mathematics
```

---

## Quick Start

### 1. Prerequisites
- Go 1.21+
- PostgreSQL (or free [Neon.tech](https://neon.tech) cloud DB)

### 2. Clone and install
```bash
git clone https://github.com/HarperKollins/AfriLearn.git
cd eduscrape
go mod tidy
```

### 3. Configure environment
```bash
cp .env.example .env
# Edit .env with your database connection string (DB_URL)
```

### 4. Deploy database schema
```bash
go run cmd/migrate/main.go
```

### 5. Ingest curriculum datasets
```bash
go run cmd/seeder/main.go
```

### 6. Start the API server
```bash
go run cmd/api/main.go
```

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
| GET | `/api/v1/` | API info and endpoint index |
| GET | `/api/v1/subjects` | List all 46 subjects, university degrees & polytechnic diploma programs |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all 22 exam boards, polytechnics & universities |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum tree with topics, subtopics & objectives |
| GET | `/api/v1/curriculum/:board/:subject/llm-prompt` | **🤖 AI Tutor Endpoint**: Formatted System Prompt, Context Window, and Module Blocks for GPT-4, Claude, Gemini, and LLaMA |
| GET | `/api/v1/search?q=:query` | Search topics across all exam boards, university degrees, and polytechnic diplomas |

---

## AI Tutor LLM Response Schema (`/llm-prompt`)

```json
{
  "success": true,
  "data": {
    "exam_board": "WAEC",
    "exam_board_slug": "waec",
    "subject": "Physics",
    "subject_slug": "physics",
    "level": "senior-secondary",
    "system_prompt": "You are an expert AI Tutor specialized in the official WAEC Physics curriculum for Senior Secondary School...",
    "topics_summary": "1. Interaction of Matter, Space, and Time\n2. Energy\n...",
    "full_context_window": "# WAEC Physics Official Curriculum Context\n\n## System Directive...\n",
    "formatted_modules": [
      {
        "module_name": "Interaction of Matter, Space, and Time",
        "difficulty": "medium",
        "llm_instruction": "Teach 'Interaction of Matter, Space, and Time' with focus on: Concepts of Matter...",
        "subtopics": ["Concepts of Matter and Position", "Units and Kinematics"]
      }
    ]
  }
}
```

---

## Technical Architecture

- **Go (Golang)** — High-performance, low-latency compiled runtime
- **Gin** — Fast HTTP framework with CORS & middleware
- **PostgreSQL (Neon)** — Relational storage with JSON array batching optimization (`pq.Array`)
- **Scraper Engine** — Modular `Scraper` interface pattern (`internal/scraper`)
