package main

import (
	"fmt"
	"strings"
)

type FunctionValues struct {
	xValues []float64
	yValues []float64
}

func (fv *FunctionValues) AddValuePair(x, y float64) {
	fv.xValues = append(fv.xValues, x)
	fv.yValues = append(fv.yValues, y)
}

func (fv *FunctionValues) Len() int {
	return len(fv.yValues)
}

func (fv *FunctionValues) String() string {
	sb := strings.Builder{}

	for i := 0; i < fv.Len(); i++ {
		buffer := fmt.Sprintf("| %7.3f | %7.3f |\n", fv.xValues[i], fv.yValues[i])
		sb.WriteString(buffer)
	}

	return sb.String()
}

func (fv *FunctionValues) GetY(x float64) float64 {
	bottom, top := 0, len(fv.xValues)-1
	var middle int

	for bottom < top {
		middle = (bottom + top) / 2
		if x > fv.xValues[middle] {
			bottom = middle + 1
		} else {
			top = middle
		}
	}

	return fv.yValues[middle]
}

func (fv FunctionValues) MultiplyBy(other FunctionValues) FunctionValues {
	for i := range fv.yValues {
		fv.yValues[i] *= other.yValues[i]
	}

	return fv
}
