package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// LLMPromptResponse represents the formatted response tailored for LLM / AI Tutor ingestion
type LLMPromptResponse struct {
	ExamBoard               string            `json:"exam_board"`
	ExamBoardSlug           string            `json:"exam_board_slug"`
	Subject                 string            `json:"subject"`
	SubjectSlug             string            `json:"subject_slug"`
	Level                   string            `json:"level"`
	SystemPrompt            string            `json:"system_prompt"`
	TopicsSummary           string            `json:"topics_summary"`
	FullContextWindow       string            `json:"full_context_window"`
	EstimatedTokenCount     int               `json:"estimated_token_count"`
	FormattedModules        []LLMModuleBlock  `json:"formatted_modules"`
	BloomsTaxonomyBreakdown map[string]int    `json:"blooms_taxonomy_breakdown"`
	DifficultyProgression   []string          `json:"difficulty_progression"`
	PedagogicalDirectives   []string          `json:"pedagogical_directives"`
	SubjectSpecificRules    []string          `json:"subject_specific_rules"`
	AdaptiveLearningPath    []string          `json:"adaptive_learning_path"`
	SuggestedChunkingNote   string            `json:"suggested_chunking_note"`
}

type LLMModuleBlock struct {
	ModuleName     string   `json:"module_name"`
	Difficulty     string   `json:"difficulty"`
	LLMInstruction string   `json:"llm_instruction"`
	Subtopics      []string `json:"subtopics"`
}

// GetLLMPrompt formats full curriculum into LLM System Prompt & Context Window
// GET /api/v1/curriculum/:board/:subject/llm-prompt
func GetLLMPrompt(c *gin.Context) {
	boardSlug := c.Param("board")
	subjectSlug := c.Param("subject")

	cacheKey := fmt.Sprintf("prompt:%s:%s", boardSlug, subjectSlug)
	if cachedVal, found := cache.GetCache().Get(cacheKey); found {
		if resp, ok := cachedVal.(models.APIResponse); ok {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	// 1. Fetch curriculum metadata
	var curr models.Curriculum
	var board models.ExamBoard
	var subject models.Subject

	err := database.DB.QueryRow(`
		SELECT 
			c.id, c.exam_board_id, c.subject_id, c.year, c.level, c.source_url, c.created_at, c.updated_at,
			eb.slug, eb.name, eb.full_name, eb.country, eb.description, eb.website,
			s.slug, s.name, s.description, s.category
		FROM curricula c
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2
		ORDER BY c.year DESC
		LIMIT 1
	`, boardSlug, subjectSlug).Scan(
		&curr.ID, &curr.ExamBoardID, &curr.SubjectID, &curr.Year, &curr.Level, &curr.SourceURL, &curr.CreatedAt, &curr.UpdatedAt,
		&board.Slug, &board.Name, &board.FullName, &board.Country, &board.Description, &board.Website,
		&subject.Slug, &subject.Name, &subject.Description, &subject.Category,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("Curriculum not found for %s/%s", boardSlug, subjectSlug),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch curriculum for LLM prompt generation",
		})
		return
	}

	// 2. Fetch topics
	topicRows, err := database.DB.Query(`
		SELECT id, curriculum_id, slug, name, description, order_index, difficulty, created_at, updated_at
		FROM topics WHERE curriculum_id = $1 ORDER BY order_index
	`, curr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch curriculum topics",
		})
		return
	}
	defer topicRows.Close()

	var topics []models.Topic
	var topicIDs []string
	for topicRows.Next() {
		var t models.Topic
		if err := topicRows.Scan(&t.ID, &t.CurriculumID, &t.Slug, &t.Name, &t.Description, &t.OrderIndex, &t.Difficulty, &t.CreatedAt, &t.UpdatedAt); err == nil {
			topics = append(topics, t)
			topicIDs = append(topicIDs, t.ID)
		}
	}

	// 3. Batch query subtopics
	subtopicMap := make(map[string]*models.Subtopic)
	var subtopicIDs []string
	if len(topicIDs) > 0 {
		subRows, err := database.DB.Query(`
			SELECT id, topic_id, slug, name, description, order_index, created_at, updated_at
			FROM subtopics WHERE topic_id = ANY($1) ORDER BY topic_id, order_index
		`, pq.Array(topicIDs))
		if err == nil {
			defer subRows.Close()
			for subRows.Next() {
				var st models.Subtopic
				if err := subRows.Scan(&st.ID, &st.TopicID, &st.Slug, &st.Name, &st.Description, &st.OrderIndex, &st.CreatedAt, &st.UpdatedAt); err == nil {
					subtopicIDs = append(subtopicIDs, st.ID)
					subtopicMap[st.ID] = &st
				}
			}
		}
	}

	// 4. Batch query objectives
	if len(subtopicIDs) > 0 {
		objRows, err := database.DB.Query(`
			SELECT id, subtopic_id, description, verb, order_index, created_at
			FROM learning_objectives WHERE subtopic_id = ANY($1) ORDER BY subtopic_id, order_index
		`, pq.Array(subtopicIDs))
		if err == nil {
			defer objRows.Close()
			for objRows.Next() {
				var obj models.LearningObjective
				if err := objRows.Scan(&obj.ID, &obj.SubtopicID, &obj.Description, &obj.Verb, &obj.OrderIndex, &obj.CreatedAt); err == nil {
					if st, exists := subtopicMap[obj.SubtopicID]; exists {
						st.Objectives = append(st.Objectives, obj)
					}
				}
			}
		}
	}

	// Group subtopics into topics in strict order_index order
	topicSubtopicMap := make(map[string][]models.Subtopic)
	for _, subID := range subtopicIDs {
		if stPtr, ok := subtopicMap[subID]; ok {
			topicSubtopicMap[stPtr.TopicID] = append(topicSubtopicMap[stPtr.TopicID], *stPtr)
		}
	}

	for i := range topics {
		if subs, ok := topicSubtopicMap[topics[i].ID]; ok {
			topics[i].Subtopics = subs
		} else {
			topics[i].Subtopics = []models.Subtopic{}
		}
	}

	// Calculate Bloom's taxonomy breakdown
	bloomsBreakdown := map[string]int{
		"remember":   0,
		"understand": 0,
		"apply":      0,
		"analyze":    0,
		"evaluate":   0,
		"create":     0,
	}

	for _, st := range subtopicMap {
		for _, obj := range st.Objectives {
			level := mapVerbToBloomLevel(obj.Verb)
			bloomsBreakdown[level]++
		}
	}

	// Level-specific instruction rules & directives
	var levelRule string
	var directives []string
	switch strings.ToLower(board.Slug) {
	case "bece", "nerdc":
		levelRule = "Strictly use simple, beginner-friendly explanations appropriate for Junior Secondary (JSS1-JSS3) students. Use simple word equations (e.g. Carbon Dioxide + Water -> Glucose + Oxygen) instead of balanced chemical formulas unless explicitly requested. Focus on foundational concepts without overloading with senior secondary or university details."
		directives = []string{
			"Use relatable everyday Nigerian analogies (e.g., local markets, solar lights, farm processes, river current for electricity).",
			"Break complex multi-step problems into numbered 2–3 step guides with worked examples.",
			"Avoid university or advanced secondary jargon — define any technical term the moment it appears.",
			"Confirm student understanding with simple check-in questions after each concept.",
			"Praise correct answers and gently redirect errors without discouraging the student.",
		}
	case "waec", "neco", "jamb":
		levelRule = "Provide comprehensive Senior Secondary (SS1-SS3 / UTME) depth aligned with WAEC/NECO marking schemes. Use standard scientific notation, balanced chemical equations, and exam-style practice questions."
		directives = []string{
			"Align answers with standard WAEC/JAMB marking scheme key phrase requirements — students lose marks for missing key words.",
			"Provide step-by-step mathematical working with explicit unit inclusions (marks are awarded per step).",
			"Include past-question style practice drills after explaining each key concept.",
			"For essay questions: follow the PEEL structure (Point, Evidence, Explanation, Link).",
			"For calculation errors: identify the specific step where the student went wrong, not just the final answer.",
		}
	default:
		levelRule = "Provide advanced university-level / polytechnic-level depth matching the official NUC/NBTE degree benchmarks."
		directives = []string{
			"Emphasize academic rigor — use correct technical terminology and standard notation throughout.",
			"Encourage critical analysis: compare and contrast theoretical models, not just description.",
			"Connect concepts to real African industry, governance, or healthcare contexts to improve retention.",
			"Encourage independent problem solving — guide with Socratic questions before giving full solutions.",
			"Uphold professional ethics relevant to the discipline in all case discussions.",
		}
	}

	// Subject-specific deep pedagogical rules
	subjectSpecificRules := buildSubjectSpecificRules(subject.Slug, board.Slug)

	systemPrompt := fmt.Sprintf(
		"You are an expert AI Tutor specialized in the official %s (%s) %s curriculum (%s level). "+
			"Your primary instruction is to explain concepts, solve practice problems, and guide students strictly aligned with "+
			"the official %s syllabus standards. %s Always provide clear, step-by-step explanations with relevant African examples.",
		board.Name, board.FullName, subject.Name, curr.Level, board.Name, levelRule,
	)

	var topicsSummaryBuilder strings.Builder
	var fullContextBuilder strings.Builder
	var moduleBlocks []LLMModuleBlock
	var adaptivePath []string
	var difficultyProgression []string

	fullContextBuilder.WriteString(fmt.Sprintf("# %s — %s Official Curriculum Context\n\n", board.Name, subject.Name))
	fullContextBuilder.WriteString(fmt.Sprintf("**Level**: %s | **Category**: %s | **Source**: %s\n\n", curr.Level, subject.Category, curr.SourceURL))
	fullContextBuilder.WriteString("## System Directive for AI Tutor\n")
	fullContextBuilder.WriteString(systemPrompt + "\n\n")
	if len(subjectSpecificRules) > 0 {
		fullContextBuilder.WriteString("## Subject-Specific Teaching Rules\n")
		for _, rule := range subjectSpecificRules {
			fullContextBuilder.WriteString(fmt.Sprintf("- %s\n", rule))
		}
		fullContextBuilder.WriteString("\n")
	}
	fullContextBuilder.WriteString("## Complete Syllabus Breakdown & Learning Objectives\n\n")

	for i, t := range topics {
		topicsSummaryBuilder.WriteString(fmt.Sprintf("%d. %s (%d subtopics, %s)\n", i+1, t.Name, len(t.Subtopics), t.Difficulty))
		adaptivePath = append(adaptivePath, fmt.Sprintf("Step %d: Master %s [%s]", i+1, t.Name, strings.ToUpper(t.Difficulty)))
		difficultyProgression = append(difficultyProgression, strings.ToUpper(t.Difficulty))

		fullContextBuilder.WriteString(fmt.Sprintf("### Module %d: %s\n", i+1, t.Name))
		fullContextBuilder.WriteString(fmt.Sprintf("*Difficulty*: %s\n", t.Difficulty))
		if t.Description != "" {
			fullContextBuilder.WriteString(fmt.Sprintf("*Description*: %s\n\n", t.Description))
		}

		var subtopicNames []string
		for j, st := range t.Subtopics {
			subtopicNames = append(subtopicNames, st.Name)
			fullContextBuilder.WriteString(fmt.Sprintf("#### Unit %d.%d: %s\n", i+1, j+1, st.Name))
			if len(st.Objectives) > 0 {
				fullContextBuilder.WriteString("Learning Objectives:\n")
				for _, obj := range st.Objectives {
					fullContextBuilder.WriteString(fmt.Sprintf("- [%s] %s\n", strings.ToUpper(obj.Verb), obj.Description))
				}
			}
			fullContextBuilder.WriteString("\n")
		}

		moduleBlocks = append(moduleBlocks, LLMModuleBlock{
			ModuleName:     t.Name,
			Difficulty:     t.Difficulty,
			LLMInstruction: fmt.Sprintf("Teach '%s' with focus on: %s", t.Name, strings.Join(subtopicNames, ", ")),
			Subtopics:      subtopicNames,
		})
	}

	fullContext := fullContextBuilder.String()
	// Rough token estimate: ~4 characters per token for English educational text
	tokenEstimate := len(fullContext) / 4

	// Chunking note based on token count
	chunkingNote := ""
	switch {
	case tokenEstimate < 2000:
		chunkingNote = fmt.Sprintf("~%d tokens — fits in a single LLM context call. Use full_context_window directly.", tokenEstimate)
	case tokenEstimate < 8000:
		chunkingNote = fmt.Sprintf("~%d tokens — fits in most 8K context models (GPT-3.5, Gemini Flash). Use full_context_window directly or split by module.", tokenEstimate)
	case tokenEstimate < 32000:
		chunkingNote = fmt.Sprintf("~%d tokens — use GPT-4o (128K), Gemini 1.5 Pro (1M), or Claude 3 for full context. Otherwise split by topic module (1 module = 1 RAG chunk).", tokenEstimate)
	default:
		chunkingNote = fmt.Sprintf("~%d tokens — too large for single-call injection. Use /embeddings endpoint to get per-module chunks, embed with OpenAI/Gemini, and retrieve relevant chunks via RAG at query time.", tokenEstimate)
	}

	promptData := LLMPromptResponse{
		ExamBoard:               board.Name,
		ExamBoardSlug:           board.Slug,
		Subject:                 subject.Name,
		SubjectSlug:             subject.Slug,
		Level:                   curr.Level,
		SystemPrompt:            systemPrompt,
		TopicsSummary:           topicsSummaryBuilder.String(),
		FullContextWindow:       fullContext,
		EstimatedTokenCount:     tokenEstimate,
		FormattedModules:        moduleBlocks,
		BloomsTaxonomyBreakdown: bloomsBreakdown,
		DifficultyProgression:   difficultyProgression,
		PedagogicalDirectives:   directives,
		SubjectSpecificRules:    subjectSpecificRules,
		AdaptiveLearningPath:    adaptivePath,
		SuggestedChunkingNote:   chunkingNote,
	}

	apiResp := models.APIResponse{
		Success: true,
		Data:    promptData,
		Meta: &models.Meta{
			Source:  curr.SourceURL,
			Version: "v1",
		},
	}

	cache.GetCache().Set(cacheKey, apiResp, 0)
	c.JSON(http.StatusOK, apiResp)
}

func mapVerbToBloomLevel(verb string) string {
	lower := strings.ToLower(strings.TrimSpace(verb))
	switch lower {
	case "define", "identify", "state", "list", "name", "recall":
		return "remember"
	case "explain", "describe", "interpret", "distinguish", "compare", "relate", "illustrate", "demonstrate":
		return "understand"
	case "apply", "calculate", "solve", "perform", "use", "express", "find", "draw":
		return "apply"
	case "analyze", "analyse", "differentiate":
		return "analyze"
	case "evaluate", "determine":
		return "evaluate"
	case "create", "construct", "formulate", "design", "build":
		return "create"
	default:
		return "understand"
	}
}

// buildSubjectSpecificRules returns deep discipline-specific pedagogical rules for the AI tutor.
// These are the substantive, non-generic rules that make AI tutoring actually useful for
// domain-specific subjects at university and polytechnic level.
func buildSubjectSpecificRules(subjectSlug, boardSlug string) []string {
	switch subjectSlug {
	case "law", "llb-bachelor-of-laws":
		return []string{
			"Always cite the full case name, court, and year for every case law reference (e.g., Donoghue v Stevenson [1932] AC 562 (HL)).",
			"Distinguish between ratio decidendi (binding legal reason) and obiter dicta (non-binding remarks) when discussing case law.",
			"Apply the IRAC method for problem questions: Issue → Rule → Application → Conclusion.",
			"Reference the Nigerian Constitution 1999 (as amended), relevant statutes (e.g., CAMA 2020, Criminal Code Act), and subsidiary legislation by chapter and section number.",
			"Note jurisdictional differences between English common law and Nigerian statutory modifications.",
			"For Equity topics: always trace the development from Court of Chancery through to modern Nigerian equity principles.",
			"For Criminal Law: always state the actus reus (guilty act) and mens rea (guilty mind) elements before discussing liability.",
			"For Land Law: distinguish between legal and equitable interests, and reference the Land Use Act 1978 (Nigeria) explicitly.",
		}

	case "medicine-surgery", "mbbs":
		return []string{
			"Present clinical cases using the SOAP note format: Subjective (patient history), Objective (examination findings), Assessment (differential diagnosis), Plan (investigations + management).",
			"Always state the drug name (generic), mechanism of action, indication, contraindications, and common adverse effects when discussing pharmacology.",
			"For differential diagnoses: start with the most life-threatening condition first (rule of worst-first).",
			"Use SI units for all laboratory values (e.g., haemoglobin in g/dL, glucose in mmol/L, creatinine in μmol/L).",
			"Reference Nigerian disease epidemiology (malaria, typhoid, sickle cell, tuberculosis) as primary context before introducing rarer global conditions.",
			"For surgical topics: describe patient positioning, landmarks, incision, and potential complications in sequence.",
			"Always include patient safety principles (e.g., WHO Surgical Safety Checklist, hand hygiene protocol) when relevant.",
			"For paediatric topics: always adjust drug dosages by weight (mg/kg) and note age-specific normal ranges.",
		}

	case "nursing-science":
		return []string{
			"Apply the Nursing Process (ADPIE) framework: Assessment → Diagnosis → Planning → Implementation → Evaluation.",
			"Always calculate drug dosages using the formula: Dose required / Stock dose × Volume.",
			"State the '5 Rights' of medication administration: Right Patient, Right Drug, Right Dose, Right Route, Right Time.",
			"For patient education: use plain language (health literacy level), and check understanding with teach-back method.",
			"Reference WHO and FMOH (Federal Ministry of Health Nigeria) guidelines for infection control protocols.",
			"For vital signs: always include normal ranges and when to escalate (MEWS/NEWS criteria).",
			"For mental health nursing: apply therapeutic communication techniques (active listening, open-ended questions, avoiding judgment).",
		}

	case "computer-science", "software-engineering":
		return []string{
			"Write pseudocode before code when explaining algorithms — pseudocode is language-agnostic and faster to grasp.",
			"Always state time complexity (Big-O) and space complexity when discussing algorithms and data structures.",
			"For database topics: show both the conceptual ER diagram and the physical SQL schema.",
			"For OOP: always illustrate with a real-world analogy (e.g., class Car with attributes and methods) before abstract definitions.",
			"For networking: work from the OSI model layer-by-layer, mapping each protocol to its layer.",
			"For software engineering: apply SOLID principles and explain which principle each design decision follows.",
			"Show working code examples in Python or JavaScript for algorithms; these are the most accessible to Nigerian CS students.",
			"For security topics: always show both the attack vector and the mitigation technique side-by-side.",
		}

	case "mechanical-engineering", "electrical-engineering", "civil-engineering",
		"chemical-engineering", "petroleum-engineering", "computer-engineering-technology":
		return []string{
			"Always state given quantities, required quantities, and relevant formula before solving any numerical problem.",
			"Use consistent SI units throughout all calculations — note unit conversions explicitly.",
			"Include safety factors and engineering tolerances in design calculations (e.g., factor of safety ≥ 2 for structural members).",
			"For materials problems: state material properties (yield strength, Young's modulus, etc.) from standard tables.",
			"Apply free body diagrams (FBD) for all structural and mechanical equilibrium problems.",
			"Reference Nigerian standards (SON — Standards Organisation of Nigeria) and international standards (ASTM, BS, ISO) where applicable.",
			"For electrical problems: always draw the circuit diagram before writing Kirchhoff's equations.",
			"Discuss environmental impact and sustainability considerations for all design and process engineering topics.",
		}

	case "economics":
		return []string{
			"Always draw supply-demand diagrams when explaining price changes, shifts, and equilibrium.",
			"Distinguish clearly between positive economics (what is) and normative economics (what ought to be).",
			"Reference Nigerian economic data (CBN Statistical Bulletin, NBS reports) for local context examples.",
			"For macroeconomics: connect theory to Nigerian fiscal and monetary policy (CBN interest rate, FAAC allocations).",
			"Calculate elasticity values numerically and interpret: |PED| > 1 = elastic, |PED| < 1 = inelastic.",
			"For international trade: apply comparative advantage with worked numerical examples showing opportunity cost.",
		}

	case "accounting":
		return []string{
			"Follow IFRS (International Financial Reporting Standards) — note where Nigerian GAAP differs.",
			"Show full double-entry bookkeeping: every debit must have a corresponding credit.",
			"Always produce the full Trial Balance before preparing Final Accounts.",
			"For ratio analysis: state the formula, calculate the ratio, and interpret what it means for the business.",
			"Reference ICAN (Institute of Chartered Accountants of Nigeria) professional standards in audit and ethics topics.",
			"For taxation topics: use current Nigerian tax rates (CIT, VAT, PAYE) from FIRS guidelines.",
		}

	case "mass-communication":
		return []string{
			"Apply communication models (Shannon-Weaver, Lasswell, Schramm) to real Nigerian media examples.",
			"Reference established Nigerian media organizations (NTA, Channels TV, The Punch, Vanguard) for examples.",
			"For journalism topics: apply the inverted pyramid structure for news writing.",
			"Discuss the BBC editorial guidelines alongside Nigerian Broadcasting Commission (NBC) Code for broadcast ethics.",
			"For PR topics: apply the RACE model (Research → Action → Communication → Evaluation).",
			"Discuss freedom of press in the context of Section 22 of the Nigerian Constitution 1999.",
		}

	case "business-administration":
		return []string{
			"Apply management theories (Taylor, Fayol, Maslow, McGregor) with real Nigerian business context.",
			"Use Porter's Five Forces and SWOT analysis frameworks for strategic management topics.",
			"Reference CAC (Corporate Affairs Commission) regulations for business law topics.",
			"Apply BCG Matrix and Ansoff Matrix for marketing strategy discussions.",
			"For HRM topics: reference Nigeria's Labour Act and Employee Compensation Act 2010.",
		}

	case "mathematics":
		return []string{
			"Always show full working — state the formula used, substitute values, simplify step-by-step.",
			"For word problems: extract given data first, define unknowns with symbols, then solve.",
			"Graph functions clearly — label axes, mark intercepts, asymptotes, and key points.",
			"For proofs: state what is to be proved, state assumptions, and conclude clearly (QED).",
			"Check answers by substituting back into the original equation.",
		}

	case "physics":
		return []string{
			"State Newton's Laws, conservation principles, or relevant formulae at the start of every mechanics problem.",
			"Use standard SI units and show all unit conversions explicitly.",
			"Draw free body diagrams for all force problems before resolving components.",
			"For wave and optics problems: sketch the wave or ray diagram as the first step.",
			"Explain the physical intuition behind each formula — what does it mean physically, not just mathematically.",
		}

	case "chemistry":
		return []string{
			"Balance all chemical equations — verify by counting atoms on both sides.",
			"State the oxidation states of all elements in redox reactions.",
			"For organic chemistry: name compounds using IUPAC nomenclature and draw structural formulae.",
			"For calculations: show molar mass calculations from the periodic table before applying stoichiometry.",
			"Note safety precautions for all laboratory procedures (corrosive, flammable, toxic reagents).",
		}

	case "biology":
		return []string{
			"Use proper binomial nomenclature (genus species, italicised) for all organisms.",
			"Draw and label biological diagrams — unlabelled diagrams score zero in WAEC.",
			"For genetics problems: always construct the Punnett square showing all gamete combinations.",
			"Connect each biological system (circulatory, respiratory, etc.) to its homeostatic function.",
			"For ecology: define the niche, habitat, and trophic level of organisms in Nigerian ecosystems.",
		}
	}

	// No subject-specific rules for this subject — return empty (generic board-level directives apply)
	return nil
}

