# Contributing to AfriLearn Curriculum Datasets

Thank you for helping build Africa's open curriculum data primitive! 🌍

AfriLearn's mission is to become the **open-source foundational data layer** for edtech applications, AI tutors, and school portals across Africa. All raw curriculum datasets in `data/curricula/` are **100% open-source under the MIT License**.

---

## 📁 Dataset Folder Structure

All datasets live in `data/curricula/`, organized by board or institution slug:

```text
data/curricula/
├── bece/                      # Junior Secondary (JSS1-3 / NERDC standard)
│   ├── mathematics.json
│   └── basic-science.json
├── waec/                      # Senior Secondary (SS1-3 / WAEC standard)
│   ├── physics.json
│   └── chemistry.json
├── jamb/                      # UTME Entrance Syllabus
├── nuc/                       # National Universities Commission (CCMAS Degree benchmarks)
│   ├── computer-science.json
│   └── law.json
├── nbte/                      # National Board for Technical Education (Polytechnics)
├── unilag/                    # University-specific (e.g. UNILAG Law, Pharmacy)
├── futo/                      # Federal University of Technology Owerri
└── yabatech/                  # Yaba College of Technology
```

---

## 📋 Standard Curriculum JSON Template

Copy this template when adding a new subject dataset:

```json
{
  "board": "unilag",
  "subject": "pharmacy",
  "level": "tertiary-degree",
  "source_url": "https://unilag.edu.ng/academics/faculty-of-pharmacy",
  "topics": [
    {
      "name": "Year 1 (100 Level): Pharmaceutical Chemistry & Biochemistry",
      "description": "Foundational topics covering organic reaction mechanisms, functional groups, and biomolecules in pharmacy.",
      "difficulty": "medium",
      "subtopics": [
        {
          "name": "PCH 101: Basic Pharmaceutical Organic Chemistry",
          "course_code": "PCH 101",
          "credit_units": "3 Units",
          "semester": "First Semester",
          "objectives": [
            "Identify functional groups and stereochemical centers in synthetic drug molecules",
            "Explain reaction mechanisms for electrophilic aromatic substitution in aromatic drug synthesis",
            "Calculate reaction yields and purity percentages from experimental synthesis data"
          ]
        }
      ]
    }
  ]
}
```

---

## 📏 Schema Validation Rules

Every dataset must pass `go run cmd/seeder/main.go --validate-only` before PR submission. The validator checks:

| Field | Requirement | Allowed Values / Examples |
|---|---|---|
| `board` | Required, string | Must match board slug (`waec`, `jamb`, `bece`, `nuc`, `unilag`, etc.) |
| `subject` | Required, string | Standard subject slug (`mathematics`, `physics`, `law`, `pharmacy`) |
| `level` | Required, string | `"junior-secondary"`, `"senior-secondary"`, `"utme"`, `"tertiary-degree"`, `"polytechnic-diploma"` |
| `source_url` | Required, valid URL | URL to official syllabus, faculty handbook, or accreditation document |
| `topics` | Non-empty array | Minimum 1 topic required per curriculum |
| `topics[].name` | Required, string | Must be descriptive (min 5 characters) |
| `topics[].difficulty` | Required, string | **Must be exactly**: `"easy"`, `"medium"`, or `"hard"` |
| `subtopics[].name` | Required, string | Must be descriptive (min 3 characters) |
| `objectives` | Array of strings | Minimum 1 objective per subtopic recommended |
| **Bloom's Action Verb** | First word of objective | **Must start with an action verb**: `calculate`, `define`, `explain`, `apply`, `analyse`, `evaluate`, `create`, `design`, `solve`, `identify`, `list`, `state`, `describe` |

---

## 🛠️ Step-by-Step Contribution & PR Workflow

### Step 1: Fork & Clone
```bash
git clone https://github.com/YOUR_USERNAME/AfriLearn.git
cd AfriLearn
```

### Step 2: Create a Branch
```bash
git checkout -b data/add-unilag-pharmacy
```

### Step 3: Add Your Dataset
Create your file under `data/curricula/[board_slug]/[subject_slug].json`.

### Step 4: Run Schema Validation
```bash
go run cmd/seeder/main.go --validate-only
```
*Expected output*: `✅ Validation passed for data/curricula/unilag/pharmacy.json`

### Step 5: Test Full Ingestion (Local DB)
```bash
go run cmd/seeder/main.go
go test -v ./...
```

### Step 6: Commit & Push
```bash
git add data/curricula/unilag/pharmacy.json
git commit -m "data: add UNILAG Pharmacy curriculum dataset"
git push origin data/add-unilag-pharmacy
```

### Step 7: Open Pull Request
Open a PR on GitHub referencing the source institution handbook or syllabus URL. Maintainers will review within 48 hours.

---

## ⚖️ License

By contributing, you agree that your data submissions will be licensed under the **MIT License**. All curriculum datasets are 100% open data for public benefit.
