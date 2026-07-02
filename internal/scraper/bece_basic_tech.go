package scraper

// BECEBasicTechScraper implements Scraper for BECE Basic Technology (JSS1 - JSS3)
type BECEBasicTechScraper struct{}

func NewBECEBasicTechScraper() *BECEBasicTechScraper {
	return &BECEBasicTechScraper{}
}

func (s *BECEBasicTechScraper) BoardSlug() string   { return "bece" }
func (s *BECEBasicTechScraper) SubjectSlug() string { return "basic-technology" }
func (s *BECEBasicTechScraper) Level() string       { return "junior-secondary" }
func (s *BECEBasicTechScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/basic-technology"
}

func (s *BECEBasicTechScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Technology and Engineering Materials (JSS1 - JSS3)",
			Description: "Types, properties, processing, and uses of wood, metals, plastics, ceramics, and rubber.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Wood Processing and Properties", Objectives: []string{
					"Classify woods into hardwoods and softwoods with examples",
					"Describe timber felling, conversion (plain/quarter sawing), seasoning, and preservation",
				}},
				{Name: "Metals, Plastics, Ceramics, and Rubber", Objectives: []string{
					"Distinguish ferrous and non-ferrous metals and alloys (steel, brass, bronze)",
					"Describe properties and uses of thermoplastics, thermosetting plastics, ceramics, and rubber",
				}},
			},
		},
		{
			Name:        "Technical Drawing and Graphics (JSS1 - JSS3)",
			Description: "Drawing instruments, board practice, geometric constructions, angles, polygons, and orthographic projection.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Drawing Instruments and Board Practice", Objectives: []string{
					"Use drawing instruments: T-square, set-squares, compasses, dividers, and protractor",
					"Identify standard drawing paper sizes (A0, A1, A2, A3, A4) and line types",
				}},
				{Name: "Geometric Construction", Objectives: []string{
					"Construct perpendicular and parallel lines, bisect lines and angles",
					"Construct regular polygons (triangles, squares, pentagons, hexagons) and circles",
				}},
				{Name: "Pictorial and Orthographic Projection", Objectives: []string{
					"Construct isometric drawings and oblique drawings of simple blocks",
					"Draw first-angle and third-angle orthographic projections of simple objects",
				}},
			},
		},
		{
			Name:        "Tools, Machines, and Workshop Safety (JSS1 - JSS3)",
			Description: "Hand tools, woodwork/metalwork tools, workshop safety rules, maintenance, and power transmission.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Bench Tools and Workshop Safety", Objectives: []string{
					"Identify measuring, marking, cutting, holding, and driving tools in woodwork and metalwork",
					"Demonstrate workshop safety regulations, first aid, and fire prevention rules",
				}},
				{Name: "Mechanical Motion and Power Transmission", Objectives: []string{
					"Explain linear, rotary, oscillating, and reciprocating motions",
					"Describe power transmission mechanisms: belt drives, chain drives, gear drives, and hydraulics",
				}},
			},
		},
	}
}
