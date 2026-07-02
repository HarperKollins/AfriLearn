# AfriLearn Curriculum API

> The foundational data layer for African educational technology.  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB · NECO · NERDC** — structured as APIs.

---

## What This Is

An infrastructure API that exposes official Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, schools, and educational platforms.

```
GET /api/v1/curriculum/bece/mathematics
GET /api/v1/curriculum/bece/basic-science
GET /api/v1/curriculum/bece/basic-technology
GET /api/v1/curriculum/bece/english-language
GET /api/v1/curriculum/bece/social-studies
GET /api/v1/curriculum/bece/business-studies

GET /api/v1/curriculum/waec/mathematics
GET /api/v1/curriculum/waec/physics
GET /api/v1/curriculum/waec/biology
GET /api/v1/curriculum/waec/chemistry
GET /api/v1/curriculum/waec/economics
GET /api/v1/curriculum/waec/government
GET /api/v1/curriculum/waec/literature

GET /api/v1/curriculum/jamb/mathematics
GET /api/v1/curriculum/jamb/physics
GET /api/v1/curriculum/jamb/chemistry
GET /api/v1/curriculum/jamb/biology
GET /api/v1/curriculum/jamb/economics
GET /api/v1/curriculum/jamb/government

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
| GET | `/api/v1/subjects` | List all 21 subjects |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all exam boards (BECE, WAEC, JAMB, NECO, NERDC) |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum tree with topics, subtopics & objectives |
| GET | `/api/v1/search?q=:query` | Search topics across all exam boards |

---

## Currently Available Curricula (19 Live Datasets)

### Junior Secondary School (BECE / JSS1 - JSS3)
- **BECE Mathematics** (`bece/mathematics`) — 4 Themes, 11 Subtopics, 35+ Objectives (NERDC JSS)
- **BECE Basic Science** (`bece/basic-science`) — 3 Themes, 8 Subtopics, 25+ Objectives (NERDC JSS)
- **BECE Basic Technology** (`bece/basic-technology`) — 3 Themes, 7 Subtopics, 22+ Objectives (NERDC JSS)
- **BECE English Studies** (`bece/english-language`) — 4 Themes, 8 Subtopics, 25+ Objectives (NERDC JSS)
- **BECE Social Studies** (`bece/social-studies`) — 3 Themes, 6 Subtopics, 20+ Objectives (NERDC JSS)
- **BECE Business Studies** (`bece/business-studies`) — 2 Themes, 4 Subtopics, 15+ Objectives (NERDC JSS)

### Senior Secondary School (WAEC / SS1 - SS3)
- **WAEC Mathematics** (`waec/mathematics`) — 9 Sections, 40+ Subtopics, 100+ Objectives (NERDC SS)
- **WAEC Physics** (`waec/physics`) — 6 Themes, 26 Subtopics, 80+ Objectives (NERDC SS)
- **WAEC Biology** (`waec/biology`) — 4 Themes, 18 Subtopics, 60+ Objectives (NERDC SS)
- **WAEC Chemistry** (`waec/chemistry`) — 4 Themes, 11 Subtopics, 40+ Objectives (NERDC SS)
- **WAEC Economics** (`waec/economics`) — 3 Themes, 8 Subtopics, 30+ Objectives (NERDC SS)
- **WAEC Government** (`waec/government`) — 3 Themes, 6 Subtopics, 25+ Objectives (NERDC SS)
- **WAEC Literature in English** (`waec/literature`) — 3 Themes, 6 Subtopics, 20+ Objectives (NERDC SS)

### Tertiary Entry Examination (JAMB / UTME)
- **JAMB Mathematics** (`jamb/mathematics`) — 5 UTME Sections, 18 Subtopics, 60+ Objectives (JAMB IBASS)
- **JAMB Physics** (`jamb/physics`) — 4 UTME Sections, 18 Subtopics, 55+ Objectives (JAMB IBASS)
- **JAMB Chemistry** (`jamb/chemistry`) — 3 UTME Sections, 7 Subtopics, 25+ Objectives (JAMB IBASS)
- **JAMB Biology** (`jamb/biology`) — 3 UTME Sections, 7 Subtopics, 25+ Objectives (JAMB IBASS)
- **JAMB Economics** (`jamb/economics`) — 3 UTME Sections, 6 Subtopics, 20+ Objectives (JAMB IBASS)
- **JAMB Government** (`jamb/government`) — 3 UTME Sections, 6 Subtopics, 20+ Objectives (JAMB IBASS)

---

## Technical Architecture

- **Go (Golang)** — High-performance, low-latency compiled runtime
- **Gin** — Fast HTTP framework with CORS & middleware
- **PostgreSQL (Neon)** — Relational storage with JSON array batching optimization (`pq.Array`)
- **Scraper Engine** — Modular `Scraper` interface pattern (`internal/scraper`)
