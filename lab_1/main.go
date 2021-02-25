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

type Member struct {
	Coefficient float64
	Power       int
}

func (m Member) MultiplyBy(times Member) Member {
	return Member{m.Coefficient * times.Coefficient, m.Power * times.Power}
}

func (m *Member) IncreaseBy(addition Member) {
	m.Coefficient += addition.Coefficient
}

type Polynom struct {
	Members []Member
}

func (p *Polynom) GetMembers() []Member {
	return p.Members
}

func (p *Polynom) Integrate() *Polynom {
	members := p.GetMembers()

	for i := 0; i < len(members); i++ {
		members[i].Power += 1
		members[i].Coefficient /= float64(members[i].Power)
	}

	return &Polynom{members}
}

func (p *Polynom) Square() *Polynom {
	members := p.GetMembers()
	var resultingMembers []Member

	for i := 0; i < len(members); i++ {
		for j := 0; j < len(members); j++ {
			resultingMembers = append(resultingMembers, members[i].MultiplyBy(members[j]))
		}
	}

	return &Polynom{resultingMembers}
}

func (p *Polynom) Count(x float64) float64 {
	members := p.GetMembers()
	accumulator := 0.0

	for _, member := range members {
		accumulator += member.Coefficient * math.Pow(x, float64(member.Power))
	}

	return accumulator
}

func (p *Polynom) MergeWith(other *Polynom) *Polynom {
	otherMembers := other.GetMembers()
	ourMembers := p.GetMembers()

	for _, member := range otherMembers {
		existsInP := false

		for i := len(ourMembers) - 1; i >= 0 && existsInP; i-- {
			if ourMembers[i].Power == member.Power {
				ourMembers[i].IncreaseBy(member)
				existsInP = true
			}
		}

		if !existsInP {
			ourMembers = append(ourMembers, member)
		}
	}

	return &Polynom{ourMembers}
}

func PicardPolynom(power int) *Polynom {
	if power == 1 {
		onlyMember := Member{1.0 / 3, 3}
		return &Polynom{[]Member{onlyMember}}
	} else {
		previousPolynom := PicardPolynom(power - 1)
		squaredPolynom := previousPolynom.Integrate().Square()
		return previousPolynom.MergeWith(squaredPolynom)
	}
}

func Picard(to, step float64, power int) *FunctionValues {
	stepCount := (int)(math.Floor(to / step))
	polynom := PicardPolynom(power)

	xValues := make([]float64, stepCount)
	yValues := make([]float64, stepCount)

	xValues[0], yValues[0] = 0, 0

	for i := 1; i < stepCount; i++ {
		xValues[i] = float64(i) * step
		yValues[i] = polynom.Count(xValues[i])
	}

	return &FunctionValues{xValues, yValues}
}

func main() {
	step := 1e-3
	limit := 1.0
	stepCount := (int)(math.Floor(limit / step))

	euler := ForwardEuler(limit, step)
	runge := RungeKutta(limit, step)
	picard1 := Picard(limit, step, 1)
	picard2 := Picard(limit, step, 2)
	picard3 := Picard(limit, step, 3)
	picard4 := Picard(limit, step, 4)

	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("|    Шаг    |   Эйлер   |Рунге-Кутта| Пикард(1) | Пикард(2) | Пикард(3) | Пикард(4) |")
	fmt.Println("-------------------------------------------------------------------------------------")

	for i := 0; i < stepCount; i++ {
		fmt.Printf("|%11f|%11f|%11f|%11f|%11f|%11f|%11f|\n",
			euler.xValues[i],
			euler.yValues[i],
			runge.yValues[i],
			picard1.yValues[i],
			picard2.yValues[i],
			picard3.yValues[i],
			picard4.yValues[i],
		)
	}
}
