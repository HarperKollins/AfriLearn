package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

// WAECMathScraper scrapes the WAEC Mathematics syllabus
type WAECMathScraper struct {
	CurriculumID string
}

// TopicData holds raw scraped data before database insertion
type TopicData struct {
	Name        string
	Description string
	Subtopics   []SubtopicData
}

// SubtopicData holds raw subtopic data
type SubtopicData struct {
	Name       string
	Objectives []string
}

// Run executes the scraper pipeline
func (s *WAECMathScraper) Run() error {
	log.Println("🕷️  Starting WAEC Mathematics scraper...")

	// Step 1: Ensure curriculum record exists and get its ID
	currID, err := s.ensureCurriculum()
	if err != nil {
		return fmt.Errorf("failed to ensure curriculum: %w", err)
	}
	s.CurriculumID = currID
	log.Printf("✅ Curriculum ID: %s\n", currID)

	// Step 2: Use hardcoded structured syllabus data from WAEC official sources
	// This is the "Golden Record" — the verified, clean dataset
	topics := s.getWAECMathTopics()

	// Step 3: Insert all topics into the database
	for i, topic := range topics {
		topicID, err := s.insertTopic(topic, i+1)
		if err != nil {
			log.Printf("⚠️  Failed to insert topic '%s': %v\n", topic.Name, err)
			continue
		}
		log.Printf("  ✅ Topic [%d/%d]: %s\n", i+1, len(topics), topic.Name)

		for j, subtopic := range topic.Subtopics {
			subtopicID, err := s.insertSubtopic(topicID, subtopic, j+1)
			if err != nil {
				log.Printf("    ⚠️  Failed to insert subtopic '%s': %v\n", subtopic.Name, err)
				continue
			}

			for k, objective := range subtopic.Objectives {
				if err := s.insertObjective(subtopicID, objective, k+1); err != nil {
					log.Printf("      ⚠️  Failed to insert objective: %v\n", err)
				}
			}
		}
	}

	log.Println("✅ WAEC Mathematics scraping complete!")
	return nil
}

// ScrapeFromWeb uses Colly to scrape additional data from web sources
func (s *WAECMathScraper) ScrapeFromWeb(targetURL string) ([]TopicData, error) {
	var topics []TopicData
	c := colly.NewCollector(
		colly.AllowedDomains("waecsyllabus.com", "waec.org.ng", "myschool.ng"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.SetRequestTimeout(30 * time.Second)

	// Rate limiting to be a good citizen
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	var currentTopic *TopicData

	// Detect topic headers (h2, h3 elements)
	c.OnHTML("h2, h3", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if text == "" {
			return
		}
		if currentTopic != nil {
			topics = append(topics, *currentTopic)
		}
		currentTopic = &TopicData{Name: text}
	})

	// Detect subtopics from list items
	c.OnHTML("li", func(e *colly.HTMLElement) {
		if currentTopic == nil {
			return
		}
		text := strings.TrimSpace(e.Text)
		if text == "" {
			return
		}
		currentTopic.Subtopics = append(currentTopic.Subtopics, SubtopicData{
			Name:       text,
			Objectives: []string{"Understand and apply " + text},
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("⚠️  Scraper error on %s: %v\n", r.Request.URL, err)
	})

	if err := c.Visit(targetURL); err != nil {
		return nil, err
	}

	// Add the last topic
	if currentTopic != nil {
		topics = append(topics, *currentTopic)
	}

	return topics, nil
}

// ensureCurriculum creates or fetches the curriculum record for WAEC Mathematics
func (s *WAECMathScraper) ensureCurriculum() (string, error) {
	var currID string

	// Try to find existing
	err := database.DB.QueryRow(`
		SELECT c.id FROM curricula c
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects sub ON c.subject_id = sub.id
		WHERE eb.slug = 'waec' AND sub.slug = 'mathematics'
		LIMIT 1
	`).Scan(&currID)

	if err == nil {
		return currID, nil // Found existing
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	// Create new curriculum record
	var boardID, subjectID string
	if err := database.DB.QueryRow(`SELECT id FROM exam_boards WHERE slug = 'waec'`).Scan(&boardID); err != nil {
		return "", fmt.Errorf("WAEC exam board not found — run schema.sql first: %w", err)
	}
	if err := database.DB.QueryRow(`SELECT id FROM subjects WHERE slug = 'mathematics'`).Scan(&subjectID); err != nil {
		return "", fmt.Errorf("Mathematics subject not found — run schema.sql first: %w", err)
	}

	currID = uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO curricula (id, exam_board_id, subject_id, year, level, source_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, currID, boardID, subjectID, time.Now().Year(), "senior-secondary",
		"https://waecsyllabus.com/mathematics-syllabus/")

	return currID, err
}

func (s *WAECMathScraper) insertTopic(topic TopicData, order int) (string, error) {
	topicID := uuid.New().String()
	slug := slugify(topic.Name)
	difficulty := classifyDifficulty(topic.Name)

	_, err := database.DB.Exec(`
		INSERT INTO topics (id, curriculum_id, slug, name, description, order_index, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (curriculum_id, slug) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = NOW()
		RETURNING id
	`, topicID, s.CurriculumID, slug, topic.Name, topic.Description, order, difficulty)

	if err != nil {
		return "", err
	}

	// Get actual ID in case of conflict
	err = database.DB.QueryRow(`SELECT id FROM topics WHERE curriculum_id = $1 AND slug = $2`, s.CurriculumID, slug).Scan(&topicID)
	return topicID, err
}

func (s *WAECMathScraper) insertSubtopic(topicID string, subtopic SubtopicData, order int) (string, error) {
	subtopicID := uuid.New().String()
	slug := slugify(subtopic.Name)

	_, err := database.DB.Exec(`
		INSERT INTO subtopics (id, topic_id, slug, name, order_index)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (topic_id, slug) DO UPDATE SET
			name = EXCLUDED.name,
			updated_at = NOW()
	`, subtopicID, topicID, slug, subtopic.Name, order)

	if err != nil {
		return "", err
	}

	err = database.DB.QueryRow(`SELECT id FROM subtopics WHERE topic_id = $1 AND slug = $2`, topicID, slug).Scan(&subtopicID)
	return subtopicID, err
}

func (s *WAECMathScraper) insertObjective(subtopicID, description string, order int) error {
	verb := extractBloomsVerb(description)
	_, err := database.DB.Exec(`
		INSERT INTO learning_objectives (id, subtopic_id, description, verb, order_index)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New().String(), subtopicID, description, verb, order)
	return err
}

// getWAECMathTopics returns the complete WAEC Mathematics syllabus
// as a structured dataset — sourced from official WAEC documentation
func (s *WAECMathScraper) getWAECMathTopics() []TopicData {
	return []TopicData{
		{
			Name:        "Number and Numeration",
			Description: "Covers the fundamental concepts of numbers, their types, operations, and applications.",
			Subtopics: []SubtopicData{
				{Name: "Number Bases", Objectives: []string{
					"Convert numbers from one base to another",
					"Perform basic operations in different number bases",
					"Apply number bases in real-world computing contexts",
				}},
				{Name: "Modular Arithmetic", Objectives: []string{
					"Define and apply modular arithmetic",
					"Solve problems involving congruences",
				}},
				{Name: "Fractions, Decimals and Approximations", Objectives: []string{
					"Perform operations on fractions and decimals",
					"Approximate numbers to required significant figures and decimal places",
					"Apply rounding in real-world contexts",
				}},
				{Name: "Indices and Logarithms", Objectives: []string{
					"Apply laws of indices to simplify expressions",
					"Define logarithm and apply laws of logarithms",
					"Solve exponential equations using logarithms",
				}},
				{Name: "Surds", Objectives: []string{
					"Simplify surds",
					"Rationalise the denominator of surd expressions",
					"Perform operations on surds",
				}},
				{Name: "Sequence and Series", Objectives: []string{
					"Identify arithmetic and geometric progressions",
					"Find the nth term and sum of AP and GP",
					"Solve practical problems involving sequences",
				}},
				{Name: "Sets", Objectives: []string{
					"Define sets and set notation",
					"Perform union, intersection, and complement operations",
					"Solve problems using Venn diagrams",
				}},
				{Name: "Logical Reasoning", Objectives: []string{
					"Apply basic principles of logic",
					"Construct truth tables",
					"Use logical reasoning to solve problems",
				}},
				{Name: "Ratio, Proportion, Rates and Percentages", Objectives: []string{
					"Solve problems on ratios and proportions",
					"Calculate rates of change",
					"Apply percentages in financial and everyday contexts",
				}},
				{Name: "Financial Arithmetic", Objectives: []string{
					"Calculate simple and compound interest",
					"Solve problems on hire purchase, discount, and profit/loss",
					"Apply financial arithmetic to banking and investment",
				}},
				{Name: "Variation", Objectives: []string{
					"Distinguish between direct, inverse, joint, and partial variation",
					"Formulate and solve equations involving variation",
				}},
			},
		},
		{
			Name:        "Algebraic Processes",
			Description: "Covers algebraic expressions, equations, graphs, and inequalities.",
			Subtopics: []SubtopicData{
				{Name: "Algebraic Expressions", Objectives: []string{
					"Simplify algebraic expressions",
					"Expand and factorise expressions",
					"Perform operations on algebraic fractions",
				}},
				{Name: "Linear Equations", Objectives: []string{
					"Solve linear equations in one and two variables",
					"Solve word problems involving linear equations",
					"Change the subject of a formula",
				}},
				{Name: "Quadratic Equations", Objectives: []string{
					"Solve quadratic equations by factorisation",
					"Solve quadratic equations using the quadratic formula",
					"Solve quadratic equations by completing the square",
					"Interpret the roots of a quadratic equation",
				}},
				{Name: "Graphs of Linear and Quadratic Functions", Objectives: []string{
					"Draw graphs of linear functions",
					"Draw graphs of quadratic functions (parabolas)",
					"Use graphs to solve equations and interpret solutions",
				}},
				{Name: "Linear Inequalities", Objectives: []string{
					"Solve linear inequalities in one variable",
					"Represent solutions on a number line",
					"Solve and graph linear inequalities in two variables",
				}},
				{Name: "Algebraic Fractions", Objectives: []string{
					"Simplify algebraic fractions",
					"Perform operations on algebraic fractions",
					"Solve equations involving algebraic fractions",
				}},
				{Name: "Functions and Relations", Objectives: []string{
					"Define and identify functions and relations",
					"Distinguish between one-to-one, many-to-one functions",
					"Find composite and inverse functions",
				}},
				{Name: "Matrices and Determinants", Objectives: []string{
					"Define matrices and perform basic operations",
					"Find the determinant and inverse of a 2x2 matrix",
					"Solve simultaneous equations using matrices",
				}},
			},
		},
		{
			Name:        "Mensuration",
			Description: "Covers measurement of lengths, areas, and volumes of plane and solid shapes.",
			Subtopics: []SubtopicData{
				{Name: "Lengths and Perimeters", Objectives: []string{
					"Calculate perimeters of plane shapes",
					"Calculate arc length of circles and sectors",
					"Apply perimeter formulae to composite shapes",
				}},
				{Name: "Areas of Plane Shapes", Objectives: []string{
					"Calculate areas of triangles, rectangles, parallelograms, and trapeziums",
					"Calculate areas of circles and sectors",
					"Find areas of composite shapes",
				}},
				{Name: "Surface Areas and Volumes", Objectives: []string{
					"Calculate surface areas of cubes, cuboids, cylinders, cones, and spheres",
					"Calculate volumes of prisms, pyramids, cylinders, cones, and spheres",
					"Solve problems involving volume and surface area in context",
				}},
			},
		},
		{
			Name:        "Plane Geometry",
			Description: "Covers the properties of geometric shapes and their relationships.",
			Subtopics: []SubtopicData{
				{Name: "Angles and Lines", Objectives: []string{
					"Identify types of angles and their properties",
					"Apply angle properties of parallel lines",
					"Solve problems involving angles at a point",
				}},
				{Name: "Triangles and Polygons", Objectives: []string{
					"Apply the properties of triangles",
					"Prove and apply congruence and similarity of triangles",
					"Calculate the sum of interior and exterior angles of polygons",
				}},
				{Name: "Circles and Circle Theorems", Objectives: []string{
					"State and apply circle theorems",
					"Solve problems involving angles in circles",
					"Apply the tangent-radius and tangent-tangent theorems",
				}},
				{Name: "Construction", Objectives: []string{
					"Construct triangles and quadrilaterals",
					"Bisect angles and line segments",
					"Construct perpendiculars and parallels",
				}},
				{Name: "Loci", Objectives: []string{
					"Identify and describe loci of points",
					"Construct loci satisfying given conditions",
					"Solve problems involving loci and regions",
				}},
			},
		},
		{
			Name:        "Coordinate Geometry",
			Description: "Covers the geometry of lines and curves using coordinates.",
			Subtopics: []SubtopicData{
				{Name: "Coordinate Geometry of Straight Lines", Objectives: []string{
					"Calculate the distance between two points",
					"Find the midpoint of a line segment",
					"Calculate the gradient (slope) of a line",
					"Determine the equation of a straight line",
					"Determine conditions for parallelism and perpendicularity",
				}},
			},
		},
		{
			Name:        "Trigonometry",
			Description: "Covers the study of relationships between angles and sides of triangles.",
			Subtopics: []SubtopicData{
				{Name: "Trigonometric Ratios", Objectives: []string{
					"Define sine, cosine, and tangent for angles in right triangles",
					"Use trigonometric tables and calculators",
					"Apply SOHCAHTOA to solve right triangles",
				}},
				{Name: "Angles of Elevation and Depression", Objectives: []string{
					"Distinguish between angles of elevation and depression",
					"Solve practical problems involving angles of elevation and depression",
				}},
				{Name: "Bearings", Objectives: []string{
					"Express directions using compass bearings",
					"Calculate distances and directions using bearings",
					"Solve navigation problems using bearings",
				}},
				{Name: "Sine and Cosine Rules", Objectives: []string{
					"State and apply the sine rule",
					"State and apply the cosine rule",
					"Solve triangles using the sine and cosine rules",
				}},
			},
		},
		{
			Name:        "Introductory Calculus",
			Description: "Covers the basic concepts of differentiation and integration.",
			Subtopics: []SubtopicData{
				{Name: "Differentiation", Objectives: []string{
					"Find the derivative of polynomial functions",
					"Apply differentiation to find gradients of curves",
					"Determine turning points (maxima and minima) of functions",
					"Apply differentiation to real-world rate-of-change problems",
				}},
				{Name: "Integration", Objectives: []string{
					"Find the indefinite integral of polynomial functions",
					"Evaluate definite integrals",
					"Apply integration to find areas under curves",
				}},
			},
		},
		{
			Name:        "Statistics and Probability",
			Description: "Covers data collection, analysis, and the mathematics of chance.",
			Subtopics: []SubtopicData{
				{Name: "Data Presentation and Analysis", Objectives: []string{
					"Collect, tabulate, and present data",
					"Draw and interpret bar charts, pie charts, histograms, and frequency polygons",
					"Calculate measures of central tendency: mean, median, mode",
					"Calculate measures of spread: range, variance, standard deviation",
				}},
				{Name: "Cumulative Frequency", Objectives: []string{
					"Draw cumulative frequency curves (ogives)",
					"Use ogives to estimate median, quartiles, and percentiles",
					"Calculate interquartile range from an ogive",
				}},
				{Name: "Probability", Objectives: []string{
					"Define probability and sample space",
					"Calculate probability of simple and compound events",
					"Apply addition and multiplication rules of probability",
					"Solve problems involving independent and mutually exclusive events",
				}},
			},
		},
		{
			Name:        "Vectors and Transformations",
			Description: "Covers vector algebra and geometric transformations.",
			Subtopics: []SubtopicData{
				{Name: "Vectors", Objectives: []string{
					"Define vectors and scalars",
					"Perform addition, subtraction, and scalar multiplication of vectors",
					"Resolve vectors into components",
					"Apply vectors to solve geometric problems",
				}},
				{Name: "Transformations", Objectives: []string{
					"Identify and apply reflection, rotation, translation, and enlargement",
					"Find images of points and shapes under transformations",
					"Describe transformations using matrices",
				}},
			},
		},
	}
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, ".", "")
	// Remove consecutive dashes
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// classifyDifficulty assigns a difficulty based on the topic name
func classifyDifficulty(topicName string) string {
	hard := []string{"calculus", "matrices", "vectors", "transformations", "logarithms", "surds"}
	medium := []string{"trigonometry", "coordinate", "statistics", "probability", "algebra"}

	lower := strings.ToLower(topicName)
	for _, h := range hard {
		if strings.Contains(lower, h) {
			return "hard"
		}
	}
	for _, m := range medium {
		if strings.Contains(lower, m) {
			return "medium"
		}
	}
	return "easy"
}

// extractBloomsVerb extracts the first action verb from an objective
func extractBloomsVerb(objective string) string {
	blooms := []string{"define", "identify", "calculate", "solve", "apply", "analyse", "evaluate",
		"create", "construct", "interpret", "explain", "distinguish", "state", "use", "find",
		"draw", "describe", "determine", "perform", "express"}

	lower := strings.ToLower(objective)
	for _, verb := range blooms {
		if strings.HasPrefix(lower, verb) {
			return verb
		}
	}
	return "understand"
}
