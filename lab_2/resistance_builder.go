package main

import (
	"math"
)

type Resistance struct {
	Length float64
	Radius float64

	MaxTemperature float64
	Amperage       float64
}

type ResistanceBuilder struct {
	resistance *Resistance
}

func NewBasicResistanceBuilder() *ResistanceBuilder {
	return NewResistanceBuilder().
		WithLength(0.12).
		WithRadius(35 * 10e-6).
		WithMaxTemperature(2000)
}

func NewResistanceBuilder() *ResistanceBuilder {
	return &ResistanceBuilder{&Resistance{}}
}

func (b *ResistanceBuilder) WithLength(length float64) *ResistanceBuilder {
	b.resistance.Length = length
	return b
}

func (b *ResistanceBuilder) WithRadius(radius float64) *ResistanceBuilder {
	b.resistance.Radius = radius
	return b
}

func (b *ResistanceBuilder) AtAmperage(amperage float64) *ResistanceBuilder {
	b.resistance.Amperage = amperage
	return b
}

func (b *ResistanceBuilder) WithMaxTemperature(maxTemp float64) *ResistanceBuilder {
	b.resistance.MaxTemperature = maxTemp
	return b
}

func (b *ResistanceBuilder) Build() float64 {
	return b.countResistance()
}

func (b *ResistanceBuilder) countResistance() float64 {
	integralPart :=
		Integrate(func(z float64) float64 {
			temperature := b.countTemperature(z)
			conductivity := GetCondInterpolator().Interpolate(temperature)

			return z * temperature * conductivity
		})

	length := b.resistance.Length
	radius := b.resistance.Radius

	return length / (2 * math.Pi * radius * radius) / integralPart
}

func (b *ResistanceBuilder) countTemperature(z float64) float64 {
	amperage := b.resistance.Amperage
	maxTemp := b.resistance.MaxTemperature

	initialTemperature := GetTemperatureInterpolator().Interpolate(amperage)
	m := GetMInterpolator().Interpolate(amperage)

	return initialTemperature + (maxTemp-initialTemperature)*math.Pow(z, m)
}

func (b *ResistanceBuilder) Experiment(from, to, step float64) FunctionValues {
	values := FunctionValues{}
	for i := from; i <= to; i += step {
		resistance := b.AtAmperage(i).Build()
		values.AddValuePair(i, resistance)
	}

	return values
}
