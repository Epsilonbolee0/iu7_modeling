package main

const (
	pointsCond  = 3
	pointsOther = 5
)

var tempInstance *Interpolator
var mInstance *Interpolator
var condInstance *Interpolator

func GetCondInterpolator() *Interpolator {
	return getInterpolator(condInstance, "conductivity.txt", pointsCond)
}

func GetMInterpolator() *Interpolator {
	return getInterpolator(mInstance, "m.txt", pointsOther)
}

func GetTemperatureInterpolator() *Interpolator {
	return getInterpolator(tempInstance, "temperature.txt", pointsOther)
}
