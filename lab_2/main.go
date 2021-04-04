package main

func main() {
	iVal, uVal :=
		NewBasicModelBuilder().
			Used().
			From(0.0).
			To(6 * 1e-4).
			WithStep(1e-6).
			Build()

	//rVal := NewBasicResistanceBuilder().Experiment(0, 6 * 1e-4, 1e-6)

	output(iVal, "i(t)", "Amperage to time", 5)
	output(uVal, "u(t)", "Voltage to time", 2400)
	//output(rVal, "r(t)", "Resistance to time", 10)
}
