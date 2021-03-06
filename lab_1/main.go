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
	return Member{m.Coefficient * times.Coefficient, m.Power + times.Power}
}

func (m *Member) IncreaseBy(addition Member) {
	m.Coefficient += addition.Coefficient
}

type Polynom struct {
	Members []Member
}

func (p *Polynom) String() string {
	sb := strings.Builder{}

	for _, member := range p.Members[:len(p.Members)-1] {
		buffer := fmt.Sprintf("%fx^%d + ", member.Coefficient, member.Power)
		sb.WriteString(buffer)
	}

	lastMember := p.Members[len(p.Members)-1]
	buffer := fmt.Sprintf("%fx^%d\n", lastMember.Coefficient, lastMember.Power)
	sb.WriteString(buffer)

	return sb.String()
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
				break
			}
		}

		if !existsInP {
			ourMembers = append(ourMembers, member)
		}
	}

	return &Polynom{ourMembers}
}

func PicarPolynom(power int) *Polynom {
	switch power {
	case 0:
		onlyMember := Member{1.0, 2}
		return &Polynom{[]Member{onlyMember}}
	case 1:
		onlyMember := Member{1.0 / 3, 3}
		return &Polynom{[]Member{onlyMember}}
	default:
		previousPolynom := PicarPolynom(power - 1)
		squaredPolynom := previousPolynom.Square()
		return PicarPolynom(0).MergeWith(squaredPolynom).Integrate()
	}
}

func Picar(to, step float64, power int) *FunctionValues {
	stepCount := (int)(math.Floor(to / step))
	polynom := PicarPolynom(power)

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
	step := 1e-5
	limit := 0.01

	stepCount := (int)(math.Floor(limit / step))

	euler := ForwardEuler(limit, step)
	runge := RungeKutta(limit, step)
	picard1 := Picar(limit, step, 1)
	picard2 := Picar(limit, step, 2)
	picard3 := Picar(limit, step, 3)
	picard4 := Picar(limit, step, 4)

	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("|    Шаг    |   Эйлер   |Рунге-Кутта| Пикард(1) | Пикард(2) | Пикард(3) | Пикард(4) |")
	fmt.Println("-------------------------------------------------------------------------------------")

	for i := 0; i < stepCount; i++ {
		fmt.Printf("|%11.5f|%11.3f|%11.3f|%11.3f|%11.3f|%11.3f|%11.3f|\n",
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
