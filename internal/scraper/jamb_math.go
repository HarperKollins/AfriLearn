package scraper

// JAMBMathScraper implements Scraper for JAMB Mathematics
type JAMBMathScraper struct{}

func NewJAMBMathScraper() *JAMBMathScraper {
	return &JAMBMathScraper{}
}

func (s *JAMBMathScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBMathScraper) SubjectSlug() string { return "mathematics" }
func (s *JAMBMathScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBMathScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/mathematics"
}

func (s *JAMBMathScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Number and Numeration",
			Description: "JAMB UTME Section I: Number bases, fractions, percentages, indices, logarithms, surds, and set theory.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Number Bases and Operations", Objectives: []string{
					"Perform calculations in bases 2 to 10 and conversion between bases",
					"Solve equations involving unknown number bases",
				}},
				{Name: "Fractions, Decimals, Percentages, and Financial Math", Objectives: []string{
					"Perform operations on fractions, decimals, and calculate percentage errors",
					"Calculate simple and compound interest, profit and loss, ratio, rates, shares, and Value Added Tax (VAT)",
				}},
				{Name: "Indices, Logarithms, and Surds", Objectives: []string{
					"Apply laws of indices and standard form calculations",
					"Apply laws of logarithms to solve logarithmic equations",
					"Simplify surds and rationalise denominators of binomial surd expressions",
				}},
				{Name: "Sets and Venn Diagrams", Objectives: []string{
					"Perform set operations: union, intersection, complement, difference",
					"Use Venn diagrams to solve 2-set and 3-set practical problems",
				}},
			},
		},
		{
			Name:        "Algebraic Processes",
			Description: "JAMB UTME Section II: Polynomials, equations, inequalities, variation, AP/GP, and matrices.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Polynomials and Algebraic Expressions", Objectives: []string{
					"Expand, factorise, and simplify polynomial expressions",
					"Apply Remainder Theorem and Factor Theorem to factorise cubic polynomials",
				}},
				{Name: "Linear, Quadratic, and Simultaneous Equations", Objectives: []string{
					"Solve linear equations and change subject of formula",
					"Solve quadratic equations using factorisation, completing square, and quadratic formula",
					"Solve simultaneous linear-linear and linear-quadratic equations in two variables",
				}},
				{Name: "Inequalities and Functions", Objectives: []string{
					"Solve linear inequalities and quadratic inequalities in one variable",
					"Represent solutions of inequalities graphically",
					"Determine domain, range, one-to-one, composite, and inverse functions",
				}},
				{Name: "Variation, Sequences, and Series", Objectives: []string{
					"Formulate and solve direct, inverse, joint, and partial variation problems",
					"Calculate nth term and sum of Arithmetic Progression (AP) and Geometric Progression (GP)",
					"Calculate sum to infinity of a geometric series",
				}},
				{Name: "Matrices and Determinants", Objectives: []string{
					"Perform matrix addition, subtraction, and multiplication up to 3x3 matrices",
					"Calculate determinant of 2x2 and 3x3 matrices and inverse of 2x2 matrix",
					"Solve system of 2 linear equations using matrix inversion / Cramer's rule",
				}},
			},
		},
		{
			Name:        "Geometry and Trigonometry",
			Description: "JAMB UTME Section III: Plane geometry, coordinate geometry of lines/circles, trigonometry, and solid geometry.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Plane Geometry and Mensuration", Objectives: []string{
					"Calculate perimeters, areas of plane shapes, and volumes of solid shapes (prisms, cones, spheres)",
					"Apply angle theorems of triangles, polygons, and circle theorems",
				}},
				{Name: "Coordinate Geometry of Lines and Circles", Objectives: []string{
					"Calculate distance between two points, midpoint, and gradient of a line",
					"Determine equations of straight lines in gradient, intercept, and double-intercept forms",
					"Determine conditions for parallel and perpendicular lines",
					"Find equation of a circle given center (h,k) and radius r: (x-h)² + (y-k)² = r²",
				}},
				{Name: "Trigonometric Ratios and Equations", Objectives: []string{
					"Calculate trigonometric ratios for standard angles (0°, 30°, 45°, 60°, 90°, 180°, 360°)",
					"Apply sine rule and cosine rule to solve acute and obtuse triangles",
					"Solve simple trigonometric equations within [0°, 360°]",
				}},
				{Name: "Bearings and 3D Solid Geometry", Objectives: []string{
					"Calculate 3-figure compass bearings and solve navigation distance problems",
					"Calculate angles between lines and planes in 3D solids",
				}},
			},
		},
		{
			Name:        "Calculus",
			Description: "JAMB UTME Section IV: Limits, differentiation, applications of derivatives, integration, and area under curves.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Limits and Continuity", Objectives: []string{
					"Evaluate limits of algebraic functions as x approaches a finite value or infinity",
					"Identify continuous and discontinuous functions",
				}},
				{Name: "Differentiation and Applications", Objectives: []string{
					"Differentiate polynomial, trigonometric, and explicit functions using power, product, quotient, and chain rules",
					"Calculate rate of change, gradient of tangents and normals to curves",
					"Locate maximum, minimum, and inflection turning points of algebraic curves",
				}},
				{Name: "Integration and Definite Integrals", Objectives: []string{
					"Integrate polynomial and simple trigonometric functions",
					"Evaluate definite integrals and calculate area bounded by a curve and axes",
				}},
			},
		},
		{
			Name:        "Statistics and Probability",
			Description: "JAMB UTME Section V: Data representation, measures of central tendency/dispersion, permutations, combinations, and probability.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Data Presentation and Analysis", Objectives: []string{
					"Interpret bar charts, pie charts, histograms, frequency tables, and ogives",
					"Calculate mean, median, mode for grouped and ungrouped data",
					"Calculate range, mean deviation, variance, and standard deviation",
				}},
				{Name: "Permutations and Combinations", Objectives: []string{
					"Apply factorial notation n! and fundamental counting principle",
					"Calculate permutations nPr and combinations nCr in selection problems",
				}},
				{Name: "Probability", Objectives: []string{
					"Calculate probability of single, mutually exclusive, and independent events",
					"Solve conditional probability and tree diagram problems",
				}},
			},
		},
	}
}
