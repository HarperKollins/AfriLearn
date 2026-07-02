package scraper

// WAECMathScraper implements the Scraper interface for WAEC Mathematics
type WAECMathScraper struct{}

func NewWAECMathScraper() *WAECMathScraper {
	return &WAECMathScraper{}
}

func (s *WAECMathScraper) BoardSlug() string   { return "waec" }
func (s *WAECMathScraper) SubjectSlug() string { return "mathematics" }
func (s *WAECMathScraper) Level() string       { return "senior-secondary" }
func (s *WAECMathScraper) SourceURL() string {
	return "https://waecsyllabus.com/mathematics-syllabus/"
}

// Topics returns the authoritative WAEC Mathematics Golden Record
func (s *WAECMathScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Number and Numeration",
			Description: "Covers fundamental concepts of numbers, their types, operations, and applications.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Number Bases", Objectives: []string{
					"Convert numbers from one base to another",
					"Perform basic operations in different number bases",
					"Apply number bases in computing contexts",
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
					"Apply percentages in financial contexts",
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
			Difficulty:  "medium",
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
					"Draw graphs of quadratic functions",
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
					"Distinguish between one-to-one and many-to-one functions",
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
			Difficulty:  "easy",
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
			Difficulty:  "easy",
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
					"Apply tangent-radius and tangent-tangent theorems",
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
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Coordinate Geometry of Straight Lines", Objectives: []string{
					"Calculate the distance between two points",
					"Find the midpoint of a line segment",
					"Calculate the gradient of a line",
					"Determine the equation of a straight line",
					"Determine conditions for parallelism and perpendicularity",
				}},
			},
		},
		{
			Name:        "Trigonometry",
			Description: "Covers relationships between angles and sides of triangles.",
			Difficulty:  "medium",
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
					"Solve triangles using sine and cosine rules",
				}},
			},
		},
		{
			Name:        "Introductory Calculus",
			Description: "Covers basic concepts of differentiation and integration.",
			Difficulty:  "hard",
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
			Description: "Covers data collection, analysis, and probability.",
			Difficulty:  "medium",
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
			Difficulty:  "hard",
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
