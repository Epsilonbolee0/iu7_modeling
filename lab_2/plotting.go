package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
)

func output(functionValues FunctionValues, filename string, plotname string, yBorder float64) {
	p := plot.New()
	p.Title.Text = plotname
	p.Y.Scale = plot.LinearScale{}
	p.Y.Tick.Marker = plot.DefaultTicks{}
	p.X.Label.Text = filename[2 : len(filename)-1]
	p.Y.Label.Text = filename

	f := plotter.NewFunction(func(x float64) float64 {
		return functionValues.GetY(x)
	})

	f.Color = color.RGBA{R: 255, A: 255}

	p.Add(f, plotter.NewGrid())
	p.Legend.Add(filename, f)

	p.X.Min = 0
	p.X.Max = 6 * 1e-4
	p.Y.Min = -1
	p.Y.Max = yBorder

	err := p.Save(20*vg.Centimeter, 20*vg.Centimeter, "lab_2/res/"+filename+".png")
	if err != nil {
		log.Panic(err)
	}
}
