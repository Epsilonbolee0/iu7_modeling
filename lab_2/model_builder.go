package main

import "fmt"

type Model struct {
	Builder *ResistanceBuilder

	ResistanceK, InductivityK, CapacityK float64

	From, To, Step     float64
	initialI, initialU float64
}

func NewBasicModelBuilder() *ModelParamsBuilder {
	return NewModelBuilder().
		With().
		Resistor(0.25).
		Capacitor(268 * 10e-6).
		Inductor(187 * 10e-6).
		VacuumTube(NewBasicResistanceBuilder()).
		Starting().
		StartingI(0.5).
		StartingU(1400)
}

type ModelBuilder struct {
	model *Model
}

func NewModelBuilder() *ModelBuilder {
	return &ModelBuilder{&Model{}}
}

type ModelPartsBuilder struct {
	ModelBuilder
}

func (b *ModelBuilder) With() *ModelPartsBuilder {
	return &ModelPartsBuilder{*b}
}

func (b *ModelPartsBuilder) Resistor(resistance float64) *ModelPartsBuilder {
	b.model.ResistanceK = resistance
	return b
}

func (b *ModelPartsBuilder) Capacitor(capacity float64) *ModelPartsBuilder {
	b.model.CapacityK = capacity
	return b
}

func (b *ModelPartsBuilder) Inductor(inductivity float64) *ModelPartsBuilder {
	b.model.InductivityK = inductivity
	return b
}

func (b *ModelPartsBuilder) VacuumTube(builder *ResistanceBuilder) *ModelPartsBuilder {
	b.model.Builder = builder
	return b
}

func (b *ModelPartsBuilder) Amperage(amperage float64) *ModelPartsBuilder {
	b.model.Builder = b.model.Builder.AtAmperage(amperage)
	return b
}

type ModelExperimentBuilder struct {
	ModelBuilder
}

func (b *ModelBuilder) Used() *ModelExperimentBuilder {
	return &ModelExperimentBuilder{*b}
}

func (b *ModelExperimentBuilder) From(from float64) *ModelExperimentBuilder {
	b.model.From = from
	return b
}

func (b *ModelExperimentBuilder) To(to float64) *ModelExperimentBuilder {
	b.model.To = to
	return b
}

func (b *ModelExperimentBuilder) WithStep(step float64) *ModelExperimentBuilder {
	b.model.Step = step
	return b
}

type ModelParamsBuilder struct {
	ModelBuilder
}

func (b *ModelBuilder) Starting() *ModelParamsBuilder {
	return &ModelParamsBuilder{*b}
}

func (b *ModelParamsBuilder) StartingU(u float64) *ModelParamsBuilder {
	b.model.initialU = u
	return b
}

func (b *ModelParamsBuilder) StartingI(i float64) *ModelParamsBuilder {
	b.model.initialI = i
	return b
}

func (b *ModelBuilder) Build() (FunctionValues, FunctionValues) {
	return b.solveHDE()
}

func (b *ModelBuilder) solveHDE() (FunctionValues, FunctionValues) {
	from, to, delta := b.experimentParams()
	currentI, currentU := b.startingParams()

	var iValues, uValues FunctionValues

	for time := from; time < to; time += delta {
		fmt.Printf("%f of %f seconds passed: u = %f, i = %f\n", time, to, currentI, currentU)

		iValues.AddValuePair(time, currentI)
		uValues.AddValuePair(time, currentU)

		kCoeffs, mCoeffs := b.findCoeffs(currentI, currentU)
		currentI = b.iterate(currentI, kCoeffs)
		currentU = b.iterate(currentU, mCoeffs)
	}

	return iValues, uValues
}

func (b *ModelBuilder) iterate(current float64, coeffs []float64) float64 {
	delta := b.model.Step
	return current + delta*(coeffs[0]+2*coeffs[1]+2*coeffs[2]+coeffs[3])
}

func (b *ModelBuilder) startingParams() (float64, float64) {
	return b.model.initialI, b.model.initialU
}

func (b *ModelBuilder) experimentParams() (float64, float64, float64) {
	return b.model.From, b.model.To, b.model.Step
}

func (b *ModelBuilder) findCoeffs(i, u float64) ([]float64, []float64) {
	k := make([]float64, 4)
	m := make([]float64, 4)

	delta := b.model.Step

	k[0] = b.calculateF(i, u)
	m[0] = b.calculateG(i)

	k[1] = b.calculateF(i+delta*k[0]/2, u+delta*m[0]/2)
	m[1] = b.calculateG(i + delta*m[0]/2)

	k[2] = b.calculateF(i+delta*k[1]/2, u+delta*m[1]/2)
	m[2] = b.calculateG(i + delta*m[1]/2)

	k[3] = b.calculateF(i+delta*k[2], u+delta*m[2])
	m[3] = b.calculateG(i + delta*m[2])

	return k, m
}

func (b *ModelBuilder) calculateF(i, u float64) float64 {
	tubeResistance := b.model.Builder.
		AtAmperage(i).
		Build()

	upperPart := u - (b.model.ResistanceK+tubeResistance)*i

	return upperPart / b.model.InductivityK
}

func (b *ModelBuilder) calculateG(i float64) float64 {
	return -i / b.model.CapacityK
}
