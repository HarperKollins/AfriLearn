# AfriLearn Curriculum API

> The foundational data layer for African educational technology.  
> **WAEC · JAMB · NECO · NERDC** — structured as APIs.

---

## What This Is

An infrastructure API that exposes Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, schools, and educational platforms.

```
GET /api/v1/curriculum/waec/mathematics
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
git clone https://github.com/yourname/curriculum-api
cd curriculum-api
go mod tidy
```

### 3. Configure environment
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 4. Set up the database
```bash
# Run schema.sql against your PostgreSQL database
psql -U postgres -d afrilearn -f internal/database/schema.sql
```

### 5. Seed WAEC Mathematics data
```bash
go run cmd/seeder/main.go
```

### 6. Start the API
```bash
go run cmd/api/main.go
```

### 7. Test it
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/curriculum/waec/mathematics
curl http://localhost:8080/api/v1/search?q=quadratic
```

---

## API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | API health check |
| GET | `/api/v1/` | API info and endpoint list |
| GET | `/api/v1/subjects` | List all subjects |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all exam boards |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum with topics |
| GET | `/api/v1/search?q=:query` | Search topics |

### Example Response
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "year": 2025,
    "level": "senior-secondary",
    "exam_board": { "name": "WAEC", "full_name": "West African Examinations Council" },
    "subject": { "name": "Mathematics", "category": "science" },
    "topics": [
      {
        "name": "Algebra",
        "difficulty": "medium",
        "subtopics": [
          {
            "name": "Quadratic Equations",
            "objectives": [
              { "description": "Solve quadratic equations by factorisation", "verb": "solve" }
            ]
          }
        ]
      }
    ]
  },
  "meta": { "version": "v1" }
}
```

---

## Project Structure

```
eduscrape/
├── cmd/
│   ├── api/        → API server entry point
│   └── seeder/     → Database seeder CLI
├── internal/
│   ├── database/   → DB connection + schema.sql
│   ├── handlers/   → HTTP route handlers
│   ├── models/     → Data models
│   └── scraper/    → Web scrapers (WAEC, JAMB, NECO)
├── .env.example    → Environment template
└── README.md
```

---

## Roadmap

- [x] WAEC Mathematics — Topics, Subtopics, Learning Objectives
- [ ] WAEC Physics
- [ ] WAEC Biology  
- [ ] JAMB Mathematics
- [ ] JAMB Physics
- [ ] NERDC Primary Curriculum
- [ ] Knowledge Graph relationships (topic prerequisites)
- [ ] Semantic search (pgvector)
- [ ] API key authentication
- [ ] Developer portal

---

## Built With

- **Go** — High-performance, compiled backend
- **Gin** — Fast HTTP framework
- **PostgreSQL** — Relational data store
- **Colly** — Web scraping
