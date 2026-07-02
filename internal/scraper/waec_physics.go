package scraper

// WAECPhysicsScraper implements Scraper for WAEC Physics
type WAECPhysicsScraper struct{}

func NewWAECPhysicsScraper() *WAECPhysicsScraper {
	return &WAECPhysicsScraper{}
}

func (s *WAECPhysicsScraper) BoardSlug() string   { return "waec" }
func (s *WAECPhysicsScraper) SubjectSlug() string { return "physics" }
func (s *WAECPhysicsScraper) Level() string       { return "senior-secondary" }
func (s *WAECPhysicsScraper) SourceURL() string {
	return "https://waecsyllabus.com/physics-syllabus/"
}

func (s *WAECPhysicsScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Interaction of Matter, Space, and Time",
			Description: "Fundamental physical quantities, measurements, kinematics, dynamics, and fluid mechanics.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Fundamental and Derived Quantities", Objectives: []string{
					"Identify fundamental and derived quantities and their SI units",
					"Determine dimensions of physical quantities",
					"Use dimensional analysis to verify physical equations",
				}},
				{Name: "Position, Distance, and Displacement", Objectives: []string{
					"Distinguish between scalar and vector quantities",
					"Calculate resultant of two or more vectors using resolution",
					"Determine position coordinates and displacement vectors",
				}},
				{Name: "Motion and Equations of Motion", Objectives: []string{
					"Analyze uniform and non-uniform linear motion",
					"Derive and apply equations of uniformly accelerated motion",
					"Interpret velocity-time and displacement-time graphs",
				}},
				{Name: "Projectile Motion", Objectives: []string{
					"Explain independence of horizontal and vertical motions",
					"Calculate maximum height, time of flight, and horizontal range",
					"Apply projectile principles to practical trajectory problems",
				}},
				{Name: "Equilibrium of Forces and Friction", Objectives: []string{
					"State conditions for equilibrium of coplanar forces",
					"Calculate moment of a force and apply principle of moments",
					"Distinguish static and dynamic friction and determine coefficient of friction",
				}},
				{Name: "Fluids at Rest and Archimedes' Principle", Objectives: []string{
					"Calculate pressure in liquids and atmospheric pressure",
					"State and apply Archimedes' principle and law of flotation",
					"Determine relative density of solids and liquids using hydrometers",
				}},
			},
		},
		{
			Name:        "Energy",
			Description: "Work, power, mechanical energy, thermal physics, gas laws, and simple machines.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Work, Energy, and Power", Objectives: []string{
					"Calculate work done by constant and variable forces",
					"Apply principle of conservation of mechanical energy",
					"Calculate mechanical power and electrical power output",
				}},
				{Name: "Thermal Energy and Temperature", Objectives: []string{
					"Distinguish temperature scales and calibrate liquid-in-glass thermometers",
					"Measure specific heat capacity and latent heat using calorimetry",
					"Apply concepts of thermal expansion of solids, liquids, and gases",
				}},
				{Name: "Gas Laws", Objectives: []string{
					"State and verify Boyle's Law, Charles' Law, and Pressure Law",
					"Apply Ideal Gas Equation PV = nRT to thermodynamic calculations",
					"Explain gas behavior using Kinetic Theory of gases",
				}},
				{Name: "Heat Transfer", Objectives: []string{
					"Explain thermal conduction, convection, and radiation mechanisms",
					"Apply Leslie's cube and Dewar flask principles",
					"Calculate rate of heat transfer through composite conductors",
				}},
				{Name: "Simple Machines", Objectives: []string{
					"Define Mechanical Advantage (MA), Velocity Ratio (VR), and Efficiency",
					"Calculate VR for levers, pulleys, inclined planes, wheel and axle, and screw jack",
					"Explain energy losses due to friction in machines",
				}},
			},
		},
		{
			Name:        "Waves and Optics",
			Description: "Mechanical and electromagnetic waves, sound, light reflection, refraction, and optical instruments.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Wave Motion and Properties", Objectives: []string{
					"Distinguish transverse and longitudinal progressive waves",
					"State wave equation v = fλ and calculate frequency and period",
					"Demonstrate reflection, refraction, diffraction, and interference of waves",
				}},
				{Name: "Sound Waves", Objectives: []string{
					"Determine speed of sound in air, liquids, and solids",
					"Explain resonance and stationary waves in stretched strings and pipes",
					"Calculate fundamental frequency and harmonics in musical instruments",
				}},
				{Name: "Reflection of Light", Objectives: []string{
					"Apply laws of reflection to plane and curved mirrors",
					"Construct ray diagrams for concave and convex spherical mirrors",
					"Apply mirror formula 1/f = 1/u + 1/v to find image position and magnification",
				}},
				{Name: "Refraction of Light and Prisms", Objectives: []string{
					"Apply Snell's law and determine refractive index of glass and water",
					"Calculate critical angle and total internal reflection applications (optical fibers)",
					"Determine minimum deviation of light passing through a glass prism",
				}},
				{Name: "Lenses and Optical Instruments", Objectives: []string{
					"Construct ray diagrams for converging and diverging thin lenses",
					"Apply lens formula 1/f = 1/u + 1/v and lens power in dioptres",
					"Describe optical principles of astronomical telescope, compound microscope, and human eye",
				}},
			},
		},
		{
			Name:        "Fields and Circuits",
			Description: "Gravitational fields, electrostatics, capacitors, magnetic fields, electromagnetic induction, and AC circuits.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Gravitational Field", Objectives: []string{
					"State Newton's law of universal gravitation",
					"Calculate gravitational potential, field strength, and escape velocity",
					"Describe Kepler's laws of planetary motion and satellite orbits",
				}},
				{Name: "Electrostatics and Capacitors", Objectives: []string{
					"Apply Coulomb's Law of electric charges",
					"Calculate capacitance of parallel plate capacitors C = εA/d",
					"Calculate equivalent capacitance for series and parallel arrangements",
				}},
				{Name: "Current Electricity and DC Circuits", Objectives: []string{
					"State Ohm's Law and calculate resistance in series and parallel",
					"Apply Kirchhoff's laws to solve complex resistor networks",
					"Determine internal resistance and electromotive force (e.m.f.) of cells",
				}},
				{Name: "Magnetism and Electromagnetic Field", Objectives: []string{
					"Map magnetic fields around straight conductors, solenoids, and bar magnets",
					"Calculate magnetic force on current-carrying conductor F = BIL sin θ",
					"Apply Faraday's and Lenz's laws of electromagnetic induction",
				}},
				{Name: "Alternating Current (AC) Circuits", Objectives: []string{
					"Distinguish peak voltage/current and root-mean-square (r.m.s.) values",
					"Calculate inductive reactance XL, capacitive reactance XC, and impedance Z",
					"Explain resonance in series R-L-C circuits and transformer power efficiency",
				}},
			},
		},
		{
			Name:        "Atomic and Modern Physics",
			Description: "Atomic structure, wave-particle duality, photoelectric effect, X-rays, and nuclear physics.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Atomic Models and Energy Levels", Objectives: []string{
					"Describe Thomson, Rutherford, and Bohr atomic models",
					"Calculate photon energy E = hf and wavelength during electron transitions",
					"Explain line emission and absorption spectra",
				}},
				{Name: "Photoelectric Effect and Wave-Particle Duality", Objectives: []string{
					"State Einstein's photoelectric equation hf = W0 + Kmax",
					"Define work function, threshold frequency, and stopping potential",
					"Calculate de Broglie wavelength λ = h/p for moving particles",
				}},
				{Name: "X-Rays and Radioactivity", Objectives: []string{
					"Describe production, properties, and medical/industrial uses of X-rays",
					"Distinguish alpha, beta, and gamma radiation properties",
					"Calculate radioactive decay, half-life T1/2, and decay constant λ",
				}},
				{Name: "Nuclear Reactions and Energy", Objectives: []string{
					"Balance nuclear equations for alpha and beta decay",
					"Calculate mass defect and nuclear binding energy using E = mc²",
					"Distinguish nuclear fission and nuclear fusion in power generation",
				}},
			},
		},
		{
			Name:        "Introductory Electronics",
			Description: "Semiconductor physics, p-n junction diodes, rectification, and transistor operations.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Semiconductors and Diodes", Objectives: []string{
					"Distinguish intrinsic and extrinsic (n-type, p-type) semiconductors",
					"Explain p-n junction formation and depletion layer barrier potential",
					"Describe forward and reverse bias characteristics of semiconductor diodes",
				}},
				{Name: "Rectification and Power Supplies", Objectives: []string{
					"Construct half-wave and full-wave bridge rectifier circuits",
					"Explain smoothing capacitor filter action in DC power supply units",
				}},
				{Name: "Transistors and Integrated Circuits", Objectives: []string{
					"Describe structure and operation of bipolar junction transistors (n-p-n, p-n-p)",
					"Use transistor as an electronic switch and current amplifier",
				}},
			},
		},
	}
}
