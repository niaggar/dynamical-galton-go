package utils

import (
	"encoding/json"
	"errors"
	"os"
)

// Particles types
const (
	Peg = iota
	Particle
)

// Distributions types
const (
	PegUniformDist = iota

	PegLogarithmicDistHorizontal
	PegGaussianDistHorizontal
	PegInverseGaussianDistHorizontal
	PegSineDistHorizontal

	PegLogarithmicDistVertical
	PegGaussianDistVertical
	PegInverseGaussianDistVertical
	PegSineDistVertical

	SphericDist
	SphericGaussianDist
)

// ParticleConfig represents the configuration of the particles
type ParticleConfig struct {
	NParticles  int
	Radius      float64
	InitDeltaX  float64
	InitDeltaY  float64
	InitDeltaVx float64
	InitDeltaVy float64
}

// PegDisplacement represents the displacement of the pegs
type PegDisplacement struct {
	Displacement bool
	AmplitudeX   float64
	AmplitudeY   float64
	FrequencyX   float64
	FrequencyY   float64
}

// PegConfig represents the configuration of the pegs
type PegConfig struct {
	MinRadius    float64
	MaxRadius    float64
	Damping      float64
	Distribution int
	DeltaFactor  float64
	CenterFactor int
	Displacement PegDisplacement
}

// BoardConfig represents the configuration of the board
type BoardConfig struct {
	VerticalSpace       float64
	HorizontalSpace     float64
	NRows               int
	NCols               int
	Periodic            bool
	StartHeightParticle float64
}

// EngineConfig represents the configuration of the logic
type EngineConfig struct {
	SubSteps    int
	MaxSteps    int
	Dt          float64
	ThreadCount int
	CPUCount    int
	Gravity     [2]float64
}

// SaveConfig represents the configuration of the save
type SaveConfig struct {
	SavePaths     bool
	SaveHistogram bool
}

// Configs represents the configuration of the simulation
type Configs struct {
	ParticleConfig ParticleConfig
	PegConfig      PegConfig
	BoardConfig    BoardConfig
	EngineConfig   EngineConfig
	SaveConfig     SaveConfig
}

// LoadConfig loads the configuration from a file
func LoadConfig(route string) (*Configs, error) {
	fileName := route + "config.json"
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("error opening the configuration file")
	}

	config := Configs{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, errors.New("error decoding the configuration file")
	}

	return &config, nil
}

// CreateBaseConfig creates a new default configuration file
func CreateBaseConfig(route string) error {
	fileName := route + "config.json"
	if fileExist(fileName) {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("error creating the configuration file")
	}

	config := Configs{
		ParticleConfig: ParticleConfig{
			NParticles:  100,
			Radius:      1,
			InitDeltaX:  0.5,
			InitDeltaY:  0,
			InitDeltaVx: 1,
			InitDeltaVy: 0,
		},
		PegConfig: PegConfig{
			MinRadius:    7,
			MaxRadius:    7,
			Damping:      0.5,
			DeltaFactor:  0.1,
			CenterFactor: 0,
			Distribution: PegUniformDist,
			Displacement: PegDisplacement{
				Displacement: false,
				AmplitudeX:   0,
				AmplitudeY:   0,
				FrequencyX:   0,
				FrequencyY:   0,
			},
		},
		BoardConfig: BoardConfig{
			VerticalSpace:       20,
			HorizontalSpace:     20,
			NRows:               20,
			NCols:               25,
			Periodic:            false,
			StartHeightParticle: 20,
		},
		EngineConfig: EngineConfig{
			SubSteps:    2,
			MaxSteps:    10000,
			Dt:          0.03,
			ThreadCount: 1,
			CPUCount:    1,
			Gravity:     [2]float64{0, -9.8},
		},
		SaveConfig: SaveConfig{
			SavePaths:     true,
			SaveHistogram: true,
		},
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return errors.New("error encoding the configuration file")
	}

	return nil
}

func fileExist(name string) bool {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
