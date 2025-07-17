// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

type Bound struct {
	// Key is the quantized value of the bound.
	Key Key
	Low float64 // Low is the lower bound of the range.
}

type DDSketchBinGenerator struct {
	// Config is the configuration for the sketch.
	Config   *Config
	Bounds   []Bound
	BoundMap map[Key]*Bound
}

func NewDDSketchBinGeneratorForAgent() *DDSketchBinGenerator {
	return NewDDSketchBinGenerator(agentConfig)
}

func NewDDSketchBinGenerator(config *Config) *DDSketchBinGenerator {
	dg := &DDSketchBinGenerator{
		Config:   config,
		Bounds:   make([]Bound, 0, defaultBinListSize),
		BoundMap: make(map[Key]*Bound, defaultBinLimit),
	}
	dg.generateBounds()
	return dg
}

func (g *DDSketchBinGenerator) GetBounds() []Bound {
	return g.Bounds
}

func (g *DDSketchBinGenerator) generateBounds() {
	half := g.Config.binLimit / 2
	for i := -half; i < half; i++ {
		key := Key(i)
		low := g.Config.binLow(key)
		b := Bound{Key: key, Low: low}
		g.Bounds = append(g.Bounds, b)
		g.BoundMap[key] = &b
	}
}

func (g *DDSketchBinGenerator) GetBound(key Key) (*Bound, bool) {
	bound, ok := g.BoundMap[key]
	if !ok {
		return nil, false
	}
	return bound, true
}

func (g *DDSketchBinGenerator) GetKeyForValue(value float64) Key {
	return g.Config.key(value)
}
