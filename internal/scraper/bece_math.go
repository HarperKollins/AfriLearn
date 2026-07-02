package scraper

// BECEMathScraper implements Scraper for BECE Mathematics (JSS1 - JSS3)
type BECEMathScraper struct{}

func NewBECEMathScraper() *BECEMathScraper {
	return &BECEMathScraper{}
}

func (s *BECEMathScraper) BoardSlug() string   { return "bece" }
func (s *BECEMathScraper) SubjectSlug() string { return "mathematics" }
func (s *BECEMathScraper) Level() string       { return "junior-secondary" }
func (s *BECEMathScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/mathematics"
}

func (s *BECEMathScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Number and Numeration (JSS1 - JSS3)",
			Description: "Whole numbers, place value, prime numbers, LCM, HCF, fractions, decimals, percentages, and binary numbers.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Whole Numbers, Factors, and Multiples", Objectives: []string{
					"Identify place values of whole numbers up to billions",
					"Find Prime Factors, Highest Common Factor (HCF), and Least Common Multiple (LCM)",
					"Express numbers in standard form and scientific notation",
				}},
				{Name: "Fractions, Decimals, and Percentages", Objectives: []string{
					"Convert between proper/improper fractions, mixed numbers, decimals, and percentages",
					"Perform addition, subtraction, multiplication, and division of fractions and decimals",
					"Solve percentage profit, loss, discount, and simple interest problems",
				}},
				{Name: "Binary Numbers (Base 2)", Objectives: []string{
					"Convert numbers from base 10 to base 2 and vice versa",
					"Perform addition and subtraction of binary numbers",
				}},
			},
		},
		{
			Name:        "Basic Algebra (JSS1 - JSS3)",
			Description: "Algebraic simplification, linear equations, substitution, expansion, factorisation, and variation.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Algebraic Expressions and Substitution", Objectives: []string{
					"Translate word problems into algebraic expressions",
					"Evaluate algebraic expressions by substitution",
					"Group like terms and simplify algebraic expressions",
				}},
				{Name: "Linear Equations and Inequalities", Objectives: []string{
					"Solve linear equations in one variable",
					"Solve simultaneous linear equations using substitution and elimination",
					"Represent linear inequalities on a number line",
				}},
				{Name: "Factorisation and Quadratic Expressions", Objectives: []string{
					"Factorise algebraic expressions by taking out common factors",
					"Expand binomial expressions (x + a)(x + b)",
					"Solve simple quadratic equations of the form x² + bx + c = 0",
				}},
			},
		},
		{
			Name:        "Geometry and Mensuration (JSS1 - JSS3)",
			Description: "Plane shapes, angles, triangles, perimeter, area, 3D solids, volume, and scale drawing.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Angles and Plane Shapes", Objectives: []string{
					"Identify acute, right, obtuse, reflex, complementary, and supplementary angles",
					"Calculate angles on a straight line, parallel lines, and in triangles/quadrilaterals",
					"Apply Pythagoras' theorem to right-angled triangles",
				}},
				{Name: "Perimeter and Area of 2D Shapes", Objectives: []string{
					"Calculate perimeters and areas of rectangles, squares, triangles, parallelograms, and trapeziums",
					"Calculate circumference and area of circles and sectors",
				}},
				{Name: "Surface Area and Volume of 3D Solids", Objectives: []string{
					"Identify properties of cubes, cuboids, cylinders, triangular prisms, and cones",
					"Calculate surface area and volume of cubes, cuboids, and cylinders",
				}},
			},
		},
		{
			Name:        "Statistics and Probability (JSS1 - JSS3)",
			Description: "Data collection, frequency tables, bar charts, pie charts, mean, median, mode, and simple probability.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Data Collection and Representation", Objectives: []string{
					"Organize raw data into tally marks and frequency distribution tables",
					"Construct and interpret bar charts, pictograms, and pie charts",
				}},
				{Name: "Measures of Central Tendency", Objectives: []string{
					"Calculate mean, median, and mode for ungrouped discrete data",
					"Identify modal class and range of data sets",
				}},
				{Name: "Simple Probability", Objectives: []string{
					"Determine experimental and theoretical probability of single events (coin toss, die roll)",
					"Express probability as fractions, decimals, and percentages",
				}},
			},
		},
	}
}
