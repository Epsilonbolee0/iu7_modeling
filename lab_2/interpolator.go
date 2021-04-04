package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type Interpolator struct {
	content        FunctionValues
	pointsRequired int
}

func (inter *Interpolator) Len() int {
	return inter.content.Len()
}

func (inter *Interpolator) XValues() []float64 {
	return inter.content.xValues
}

func (inter *Interpolator) YValues() []float64 {
	return inter.content.yValues
}

func (inter *Interpolator) Values() ([]float64, []float64) {
	return inter.XValues(), inter.YValues()
}

func readInterpolationData(path string) (FunctionValues, error) {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(ex + path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	result := FunctionValues{}

	for scanner.Scan() {
		cleansedText := strings.Fields(scanner.Text())

		key, _ := strconv.ParseFloat(cleansedText[0], 64)
		value, _ := strconv.ParseFloat(cleansedText[1], 64)

		result.AddValuePair(key, value)
	}

	return result, nil
}

func getInterpolator(instance *Interpolator, filepath string, points int) *Interpolator {
	if instance == nil {
		content, e := readInterpolationData("\\lab_2\\src\\" + filepath)
		inter := Interpolator{}
		if e == nil {
			inter.content = content
			inter.pointsRequired = points
		}
		instance = &inter
	}
	return instance
}

func (inter *Interpolator) getNearest(x float64) int {
	xVal := inter.XValues()

	minDelta := math.Abs(x - xVal[0])
	minIndex := 0

	for ind, value := range xVal {
		delta := math.Abs(x - value)
		if delta < minDelta {
			minDelta = delta
			minIndex = ind
		}
	}

	return minIndex
}

func (inter *Interpolator) chooseDots(x float64) FunctionValues {
	tableLength := inter.Len()
	nearestIndex := inter.getNearest(x)
	dotsCount := inter.pointsRequired + 1

	var startIndex, endIndex int

	if dotsCount+nearestIndex+1 > tableLength {
		startIndex = tableLength - dotsCount
		endIndex = tableLength
	} else if dotsCount < nearestIndex || nearestIndex-startIndex < dotsCount {
		startIndex = 0
		endIndex = dotsCount
	} else {
		startIndex = nearestIndex - dotsCount + 1
		endIndex = startIndex + dotsCount
	}

	xVal, yVal := inter.Values()
	return FunctionValues{xVal[startIndex:endIndex], yVal[startIndex:endIndex]}
}

func (inter *Interpolator) getMatrix(fv FunctionValues) [][]float64 {
	xVal, yVal := fv.xValues, fv.yValues

	var matrix [][]float64
	matrix = append(matrix, xVal)
	matrix = append(matrix, yVal)
	length := fv.Len() - 1

	for i := 0; i < length; i++ {
		row := make([]float64, fv.Len()-i)

		for j := 0; j < length-i; j++ {
			row[j] = (matrix[i+1][j] - matrix[i+1][j+1]) / (xVal[j] - xVal[i+j+1])
		}

		matrix = append(matrix, row)
	}

	return matrix
}

func (inter *Interpolator) Interpolate(x float64) float64 {
	dots := inter.chooseDots(x)
	matrix := inter.getMatrix(dots)

	temp, accumulator := 1.0, 0.0
	for i := 0; i <= inter.pointsRequired; i++ {
		accumulator += temp * matrix[i+1][0]
		temp *= x - matrix[0][i]
	}

	return accumulator
}
