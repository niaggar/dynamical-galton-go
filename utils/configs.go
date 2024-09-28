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

// PegConfig represents the configuration of the pegs
type PegConfig struct {
	MinRadius    float64
	MaxRadius    float64
	Damping      float64
	Distribution int
	DeltaFactor  float64
	CenterFactor int
}

// BoardConfig represents the configuration of the board
type BoardConfig struct {
	VerticalSpace   float64
	HorizontalSpace float64
	NRows           int
	NCols           int
}

// EngineConfig represents the configuration of the logic
type EngineConfig struct {
	SubSteps int
	MaxSteps int
	Dt       float64
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
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("error creating the configuration file")
	}

	config := Configs{
		ParticleConfig: ParticleConfig{
			NParticles:  100,
			Radius:      5,
			InitDeltaX:  1,
			InitDeltaY:  1,
			InitDeltaVx: 5,
			InitDeltaVy: 5,
		},
		PegConfig: PegConfig{
			MinRadius:    15,
			MaxRadius:    15,
			Damping:      0.5,
			DeltaFactor:  0.1,
			CenterFactor: 0,
			Distribution: PegUniformDist,
		},
		BoardConfig: BoardConfig{
			VerticalSpace:   20,
			HorizontalSpace: 20,
			NRows:           10,
			NCols:           11,
		},
		EngineConfig: EngineConfig{
			SubSteps: 10,
			MaxSteps: 1000,
			Dt:       0.03,
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
