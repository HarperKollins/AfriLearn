# AfriLearn Curriculum API

> The foundational data layer for African educational technology.  
> **WAEC · JAMB · NECO · NERDC** — structured as APIs.

---

## What This Is

An infrastructure API that exposes official Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, schools, and educational platforms.

```
GET /api/v1/curriculum/waec/mathematics
GET /api/v1/curriculum/waec/physics
GET /api/v1/curriculum/waec/biology
GET /api/v1/curriculum/jamb/mathematics
GET /api/v1/curriculum/jamb/physics
GET /api/v1/subjects
GET /api/v1/search?q=quadratic equations
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
| GET | `/api/v1/subjects` | List all subjects |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all exam boards (WAEC, JAMB, NECO, NERDC) |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum tree with topics, subtopics & objectives |
| GET | `/api/v1/search?q=:query` | Search topics across all exam boards |

---

## Currently Available Curricula

- **WAEC Mathematics** (`waec/mathematics`) — 9 Sections, 40+ Subtopics, 100+ Objectives (NERDC Senior Secondary)
- **WAEC Physics** (`waec/physics`) — 6 Themes, 26 Subtopics, 80+ Objectives (NERDC Senior Secondary)
- **WAEC Biology** (`waec/biology`) — 4 Themes, 18 Subtopics, 60+ Objectives (NERDC Senior Secondary)
- **JAMB Mathematics** (`jamb/mathematics`) — 5 UTME Sections, 18 Subtopics, 60+ Objectives (JAMB IBASS)
- **JAMB Physics** (`jamb/physics`) — 4 UTME Sections, 18 Subtopics, 55+ Objectives (JAMB IBASS)

---

## Technical Architecture

- **Go (Golang)** — High-performance, low-latency compiled runtime
- **Gin** — Fast HTTP framework with CORS & middleware
- **PostgreSQL (Neon)** — Relational storage with JSON array batching optimization (`pq.Array`)
- **Scraper Engine** — Modular `Scraper` interface pattern (`internal/scraper`)
