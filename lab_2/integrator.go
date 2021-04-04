package main

const step = 1e-2

type function func(float64) float64

func Integrate(what function) float64 {
	return simpsonMethod(0, 1, what)
}

func simpsonMethod(from, to float64, what function) float64 {
	accumulator := what(from) + what(to-step) + what(to)
	steps := int((to-from)/step) - 3

	oddPart, evenPart := 0.0, 0.0
	doubledStep := 2 * step
	currX := 2 * step

	for i := 1; i < steps/2; i++ {
		oddPart += what(currX - step)
		evenPart -= what(currX)
		currX += doubledStep
	}

	accumulator += 4*oddPart + 2*evenPart

	return step * accumulator / 6
}
