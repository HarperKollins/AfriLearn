# Contributing to AfriLearn Curriculum Datasets

Thank you for helping democratize and standardize African educational curriculum data! 🌍

AfriLearn's mission is to become the **open-source foundational data primitive** for edtech applications, AI tutors, and school portals across Africa. All raw curriculum data in `data/curricula/` is 100% open-source under the MIT License.

---

## 📁 Repository Data Structure

All curricula are organized under `data/curricula/`:
```
data/curricula/
├── bece/                  # Junior Secondary (JSS1-3)
├── waec/                  # Senior Secondary (SS1-3)
├── jamb/                  # UTME Entrance Syllabus
├── nuc/                   # National Universities Commission (CCMAS Standards)
├── [university_slug]/     # University-specific programs (e.g., unilag, unn, futo)
└── [polytechnic_slug]/    # Polytechnic programs (e.g., yabatech, imt)
```

---

## 📜 Curriculum JSON Schema Specification

Every JSON file in `data/curricula/` must strictly follow this structure:

```json
{
  "board": "unilag",
  "subject": "law",
  "level": "tertiary-degree",
  "source_url": "https://unilag.edu.ng/academics/faculty-of-law",
  "topics": [
    {
      "name": "Year 1 (100 Level): Legal Foundations & Methodologies",
      "description": "Foundational law courses setting up legal reasoning, legal systems, and statutory interpretation.",
      "difficulty": "easy",
      "subtopics": [
        {
          "name": "Legal Methods I: Sources & Statutory Interpretation",
          "course_code": "LAW 101",
          "credit_units": "4 Units",
          "semester": "First Semester",
          "objectives": [
            "Apply common law rules of statutory interpretation including the literal, golden, and mischief rules",
            "Evaluate judicial precedent (stare decisis) and ratio decidendi versus obiter dictum in Nigerian case law"
          ]
        }
      ]
    }
  ]
}
```

### Required Fields & Rules:
1. **`board`**: Must match an existing exam board or institution slug (`waec`, `jamb`, `bece`, `nuc`, `unilag`, `futo`, etc.).
2. **`subject`**: Must match standard subject slug (`mathematics`, `physics`, `law`, `computer-science`, `medicine-and-surgery`).
3. **`difficulty`**: Must be one of `"easy"`, `"medium"`, or `"hard"`.
4. **`objectives`**: Each learning objective MUST start with an explicit Bloom's Taxonomy action verb (`calculate`, `define`, `explain`, `apply`, `analyse`, `evaluate`, `create`, `design`, `solve`).
5. **Tertiary Degree Attributes** (`course_code`, `credit_units`, `semester`): Highly encouraged for university and polytechnic courses (e.g. `"LAW 101"`, `"4 Units"`, `"First Semester"`).

---

## 🛠️ Pre-Commit Quality & Validation Check

Before submitting a Pull Request, run the local dataset validator:

```bash
# Verify schema compliance across all JSON files
go run cmd/seeder/main.go --validate-only
```

If any file contains invalid difficulty tags, missing fields, or objectives under 10 characters, the validator will highlight the exact line and file.

---

## 🚀 How to Submit a Pull Request

1. Fork the repository on GitHub: `https://github.com/HarperKollins/AfriLearn`
2. Create a new branch: `git checkout -b feature/add-unilag-engineering`
3. Add your new curriculum file under `data/curricula/[board]/[subject].json`
4. Run validation: `go run cmd/seeder/main.go --validate-only`
5. Run tests: `go test -v ./...`
6. Commit and push: `git commit -m "feat: add UNILAG Mechanical Engineering curriculum"`
7. Open a Pull Request on GitHub.
