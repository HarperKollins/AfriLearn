# AfriLearn Curriculum API

> The foundational data layer for African educational technology.  
> **BECE (JSS1-3) · WAEC (SS1-3) · JAMB · NUC CCMAS (100L-500L) · Nigerian Universities (EBSU, FUNAI, UNEC, UNN)** — structured as APIs.

---

## What This Is

An infrastructure API that exposes official Nigerian/African curriculum data as clean, structured JSON endpoints. Built for EdTech developers, AI tutor companies, universities, schools, and educational platforms.

```
# Specific University Curricula (EBSU, FUNAI, UNEC, UNN)
GET /api/v1/curriculum/ebsu/computer-science
GET /api/v1/curriculum/funai/computer-science
GET /api/v1/curriculum/unec/law
GET /api/v1/curriculum/unn/computer-science

# University Higher Education (NUC CCMAS Standards - 100L to 500L)
GET /api/v1/curriculum/nuc/computer-science
GET /api/v1/curriculum/nuc/medicine-and-surgery
GET /api/v1/curriculum/nuc/electrical-engineering
GET /api/v1/curriculum/nuc/law
GET /api/v1/curriculum/nuc/accounting
GET /api/v1/curriculum/nuc/business-administration
GET /api/v1/curriculum/nuc/nursing-science
GET /api/v1/curriculum/nuc/mechanical-engineering
GET /api/v1/curriculum/nuc/mass-communication

# Junior & Senior Secondary School
GET /api/v1/curriculum/bece/mathematics
GET /api/v1/curriculum/waec/mathematics
GET /api/v1/curriculum/jamb/mathematics

GET /api/v1/subjects
GET /api/v1/exam-boards
GET /api/v1/search?q=data structures
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
| GET | `/api/v1/subjects` | List all 43 subjects & degree programs across all 17 NUC disciplines |
| GET | `/api/v1/subjects/:slug` | Get subject by slug |
| GET | `/api/v1/exam-boards` | List all 16 exam boards & universities (BECE, WAEC, JAMB, NUC, EBSU, FUNAI, UNEC, UNN, UNILAG, UI, OAU, ABU, Covenant) |
| GET | `/api/v1/curriculum/:board/:subject` | Full curriculum tree with topics, subtopics & objectives |
| GET | `/api/v1/search?q=:query` | Search topics across all exam boards and university degrees |

---

## Currently Available Curricula (32 Live Datasets)

### Nigerian Universities (EBSU, FUNAI, UNEC, UNN)
- **Ebonyi State University, Abakaliki (EBSU)** — `ebsu/computer-science` (Faculty of Physical Sciences)
- **Alex Ekwueme Federal University, Ndufu-Alike (AE-FUNAI)** — `funai/computer-science` (Faculty of Science)
- **University of Nigeria, Enugu Campus (UNEC)** — `unec/law` (Faculty of Law, Enugu Campus)
- **University of Nigeria, Nsukka (UNN)** — `unn/computer-science` (Faculty of Physical Sciences, Nsukka)

### National Universities Commission (NUC CCMAS Standards - 100L to 500L)
- **B.Sc. Computer Science** (`nuc/computer-science`)
- **M.B.B.S. Medicine & Surgery** (`nuc/medicine-and-surgery`)
- **B.Eng. Electrical & Electronic Engineering** (`nuc/electrical-engineering`)
- **LL.B. Bachelor of Laws** (`nuc/law`)
- **B.Sc. Accounting** (`nuc/accounting`)
- **B.Sc. Business Administration** (`nuc/business-administration`)
- **B.N.Sc. Nursing Science** (`nuc/nursing-science`)
- **B.Eng. Mechanical Engineering** (`nuc/mechanical-engineering`)
- **B.Sc. Mass Communication** (`nuc/mass-communication`)

### Junior Secondary School (BECE / JSS1 - JSS3)
- **BECE Mathematics** (`bece/mathematics`)
- **BECE Basic Science** (`bece/basic-science`)
- **BECE Basic Technology** (`bece/basic-technology`)
- **BECE English Studies** (`bece/english-language`)
- **BECE Social Studies** (`bece/social-studies`)
- **BECE Business Studies** (`bece/business-studies`)

### Senior Secondary School (WAEC & NECO / SS1 - SS3)
- **WAEC Mathematics** (`waec/mathematics`)
- **WAEC Physics** (`waec/physics`)
- **WAEC Biology** (`waec/biology`)
- **WAEC Chemistry** (`waec/chemistry`)
- **WAEC Economics** (`waec/economics`)
- **WAEC Government** (`waec/government`)
- **WAEC Literature in English** (`waec/literature`)

### Tertiary Entry Examination (JAMB / UTME)
- **JAMB Mathematics** (`jamb/mathematics`)
- **JAMB Physics** (`jamb/physics`)
- **JAMB Chemistry** (`jamb/chemistry`)
- **JAMB Biology** (`jamb/biology`)
- **JAMB Economics** (`jamb/economics`)
- **JAMB Government** (`jamb/government`)

---

## Technical Architecture

- **Go (Golang)** — High-performance, low-latency compiled runtime
- **Gin** — Fast HTTP framework with CORS & middleware
- **PostgreSQL (Neon)** — Relational storage with JSON array batching optimization (`pq.Array`)
- **Scraper Engine** — Modular `Scraper` interface pattern (`internal/scraper`)
