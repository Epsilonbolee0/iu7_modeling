package main

import (
	"fmt"
	"math"
	"strings"
)

func Function(y, x float64) float64 {
	return y*y + x*x
}

type FunctionValues struct {
	xValues []float64
	yValues []float64
}

func (f *FunctionValues) String() string {
	sb := strings.Builder{}

	for i := 0; i < len(f.xValues); i++ {
		str := fmt.Sprintf(" f(%f) = %f\n", f.xValues[i], f.yValues[i])
		sb.WriteString(str)
	}

	return sb.String()
}

func ForwardEuler(to, step float64) *FunctionValues {
	stepCount := (int)(math.Floor(to / step))

	xValues := make([]float64, stepCount)
	yValues := make([]float64, stepCount)

	xValues[0], yValues[0] = 0, 0

	for i := 1; i < stepCount; i++ {
		xValues[i] = float64(i) * step
		yValues[i] = yValues[i-1] + step*Function(xValues[i-1], yValues[i-1])
	}

	return &FunctionValues{xValues, yValues}
}

func RungeKutta(to, step float64) *FunctionValues {
	stepCount := (int)(math.Floor(to / step))

	xValues := make([]float64, stepCount)
	yEuler := make([]float64, stepCount)
	yValues := make([]float64, stepCount)

	xValues[0], yValues[0] = 0, 0

	for i := 1; i < stepCount; i++ {
		xValues[i] = float64(i) * step
		yEuler[i] = yEuler[i-1] + step*Function(xValues[i-1], yEuler[i-1])
		yValues[i] = yValues[i-1] + step*Function(xValues[i-1]+step/2.0, yValues[i-1]+yEuler[i]/2)
	}

	return &FunctionValues{xValues, yValues}
}

type Polynom struct {
	Coefficient float64
	Power       int
}

func IntegratePolynom(polynom []Polynom) []Polynom {
	for i := 0; i < len(polynom); i++ {
		polynom[i].Power += 1
		polynom[i].Coefficient /= float64(polynom[i].Power)
	}

	return polynom
}

func SquarePolynom(polynom []Polynom) []Polynom {
	var result []Polynom

	for i := 0; i < len(polynom); i++ {
		for j := 0; j < len(polynom); j++ {
			result = append(result, MultiplyMembers(polynom[i], polynom[j]))
		}
	}

	return result
}

func MergePolynoms(a, b []Polynom) []Polynom {
	for _, member := range b {
		exists := false

		for i := len(a) - 1; i >= 0 && exists; i-- {
			if a[i] == member {
				a[i].Power += member.Power
				a[i].Coefficient += member.Coefficient
				exists = true
			}
		}

		if !exists {
			a = append(a, member)
		}
	}

	return a
}

func MultiplyMembers(a, b Polynom) Polynom {
	return Polynom{a.Coefficient * b.Coefficient, a.Power * b.Power}
}

func picardPolynom(power int) func() []Polynom {
	if power == 1 {
		return func() []Polynom {
			return []Polynom{{1.0 / 3, 3}}
		}
	} else {
		return func() []Polynom {
			previousPolynom := picardPolynom(power - 1)()
			squarePolynom := SquarePolynom(IntegratePolynom(previousPolynom))
			return MergePolynoms(previousPolynom, squarePolynom)
		}
	}
}

func CountPolynom(x float64, polynom []Polynom) float64 {
	accumulator := 0.0
	for _, member := range polynom {
		accumulator += member.Coefficient * math.Pow(x, float64(member.Power))
	}

	return accumulator
}

func PicardPolynom(power int) func(x float64) float64 {
	return func(x float64) float64 {
		return CountPolynom(x, picardPolynom(power)())
	}
}

func Picard(to, step float64, power int) *FunctionValues {
	stepCount := (int)(math.Floor(to / step))
	function := PicardPolynom(power)

	xValues := make([]float64, stepCount)
	yValues := make([]float64, stepCount)

	xValues[0], yValues[0] = 0, 0

	for i := 1; i < stepCount; i++ {
		xValues[i] = float64(i) * step
		yValues[i] = function(xValues[i])
	}

	return &FunctionValues{xValues, yValues}
}

func main() {
	step := 1e-3
	limit := 4.0

	fmt.Println(" Euler method: ")
	fmt.Print(ForwardEuler(limit, step).String())
	fmt.Println(" Runge-Kutta method: ")
	fmt.Print(RungeKutta(limit, step).String())
	fmt.Println(" Picard (1): ")
	fmt.Print(Picard(limit, step, 1).String())
	fmt.Println(" Picard (2): ")
	fmt.Print(Picard(limit, step, 2).String())
	fmt.Println(" Picard (3): ")
	fmt.Print(Picard(limit, step, 3).String())
	fmt.Println(" Picard (4): ")
	fmt.Print(Picard(limit, step, 4).String())
}
