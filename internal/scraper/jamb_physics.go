package scraper

// JAMBPhysicsScraper implements Scraper for JAMB Physics
type JAMBPhysicsScraper struct{}

func NewJAMBPhysicsScraper() *JAMBPhysicsScraper {
	return &JAMBPhysicsScraper{}
}

func (s *JAMBPhysicsScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBPhysicsScraper) SubjectSlug() string { return "physics" }
func (s *JAMBPhysicsScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBPhysicsScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/physics"
}

func (s *JAMBPhysicsScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Measurement and Mechanics",
			Description: "JAMB UTME Section 1: Units, vectors, kinematics, Newton's laws, gravitation, work/energy, and simple machines.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Units and Dimensions", Objectives: []string{
					"Identify SI base units and derived units",
					"Determine dimensions of physical quantities and test equation consistency",
				}},
				{Name: "Scalars, Vectors, and Kinematics", Objectives: []string{
					"Resolve vectors into rectangular components and calculate resultant",
					"Analyze linear motion, projectile motion, and uniform circular motion (centripetal force)",
				}},
				{Name: "Newton's Laws of Motion and Momentum", Objectives: []string{
					"Apply Newton's three laws of motion to real-world mechanical systems",
					"Calculate impulse and linear momentum conservation in elastic and inelastic collisions",
				}},
				{Name: "Forces, Equilibrium, and Friction", Objectives: []string{
					"Calculate gravitational force using Newton's law of universal gravitation",
					"State conditions for equilibrium of rigid bodies and calculate moments",
					"Determine static and kinetic friction coefficients on horizontal and inclined planes",
				}},
				{Name: "Work, Energy, Power, and Simple Machines", Objectives: []string{
					"Apply conservation of mechanical energy principle (PE + KE = constant)",
					"Calculate efficiency, Mechanical Advantage (MA), and Velocity Ratio (VR) of machines",
				}},
			},
		},
		{
			Name:        "Matter, Heat, and Waves",
			Description: "JAMB UTME Section 2: Fluid pressure, thermal expansion, gas laws, SHM, wave properties, sound, and geometrical optics.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Elasticity and Fluid Mechanics", Objectives: []string{
					"State Hooke's Law and calculate Young's modulus E = (F/A)/(ΔL/L)",
					"Calculate fluid pressure P = hρg and apply Pascal's and Archimedes' principles",
				}},
				{Name: "Thermal Physics and Gas Laws", Objectives: []string{
					"Convert between Celsius, Fahrenheit, and Kelvin temperature scales",
					"Measure specific heat capacity and latent heat using method of mixtures",
					"Apply Boyle's Law, Charles' Law, Pressure Law, and Ideal Gas Law PV = nRT",
				}},
				{Name: "Simple Harmonic Motion (SHM) and Waves", Objectives: []string{
					"Analyze SHM in simple pendulums and mass-spring systems T = 2π√(m/k)",
					"Apply wave equation v = fλ to transverse, longitudinal, and electromagnetic waves",
					"Explain interference, diffraction, polarization, and stationary wave resonance",
				}},
				{Name: "Sound Waves and Acoustics", Objectives: []string{
					"Determine velocity of sound in gases, liquids, and solids",
					"Calculate pitch, loudness, intensity, echo, and Doppler effect frequencies",
				}},
				{Name: "Geometrical Optics and Instruments", Objectives: []string{
					"Apply reflection and refraction laws at plane and curved spherical surfaces",
					"Calculate critical angle and total internal reflection applications",
					"Apply mirror/lens formulas 1/f = 1/u + 1/v to microscopes, telescopes, and eye defects",
				}},
			},
		},
		{
			Name:        "Electricity and Magnetism",
			Description: "JAMB UTME Section 3: Electrostatics, capacitors, DC circuits, magnetic fields, induction, and AC circuits.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Electrostatics and Capacitance", Objectives: []string{
					"Calculate electrostatic force using Coulomb's Law F = kq1q2/r²",
					"Calculate electric field intensity, potential, and energy stored in capacitors E = ½CV²",
				}},
				{Name: "Current Electricity and DC Circuits", Objectives: []string{
					"State Ohm's Law and calculate equivalent resistance in series and parallel networks",
					"Apply Kirchhoff's Current Law and Voltage Law to electrical circuits",
					"Calculate electrical energy, power P = IV = I²R, and internal resistance of cells",
				}},
				{Name: "Magnetism and Electromagnetic Induction", Objectives: []string{
					"Map magnetic fields around solenoids, straight wires, and earth's magnetic field",
					"Calculate magnetic force on moving charge F = qvB sin θ and current-carrying wire F = BIL sin θ",
					"State Faraday's and Lenz's laws and explain transformer voltage step-up/step-down Vp/Vs = Np/Ns",
				}},
				{Name: "AC Circuits", Objectives: []string{
					"Distinguish peak voltage V0 and root-mean-square Vrms = V0/√2",
					"Calculate inductive reactance XL, capacitive reactance XC, and impedance Z in RLC circuits",
				}},
			},
		},
		{
			Name:        "Atomic and Modern Physics",
			Description: "JAMB UTME Section 4: Conduction in gases, photoelectric effect, atomic models, X-rays, radioactivity, and semiconductors.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Conduction in Gases and Cathode Rays", Objectives: []string{
					"Describe discharge tube phenomena and properties of cathode rays",
					"Explain thermionic emission and cathode ray oscilloscope (CRO) applications",
				}},
				{Name: "Photoelectric Effect and Wave-Particle Duality", Objectives: []string{
					"Apply Einstein's photoelectric equation hf = hf0 + ½mv²",
					"Calculate photon energy E = hf and de Broglie wavelength λ = h/p",
				}},
				{Name: "Atomic Models and X-Rays", Objectives: []string{
					"Describe Thomson, Rutherford, and Bohr atomic spectral models",
					"Explain X-ray production in Coolidge tube, hard vs soft X-rays, and medical safety",
				}},
				{Name: "Radioactivity, Nuclear Reactions, and Electronics", Objectives: []string{
					"Calculate radioactive decay, half-life T½ = 0.693/λ, and mass defect energy E = Δmc²",
					"Distinguish intrinsic vs extrinsic p-type and n-type semiconductors and p-n junction diode rectification",
				}},
			},
		},
	}
}
