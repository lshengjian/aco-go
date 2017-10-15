package util
import (
	"image/color"
	"os"
	"fmt"
	"math/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
//	"gonum.org/v1/plot/vg/draw"
)
type CityData struct{
  Idx int
  X,Y float64 
}
type TourLengthData struct{
	Iteration int
	Data int 
  }
//prints error and exits on abnormal conditions
func PrintErrorAndExit(err error) {
	fmt.Print(err)
	os.Exit(2)
}
func SaveVisitedImage(fname string,data []CityData){
	g, err := plot.New()
	if err != nil {
		panic(err)
	}

	g.Title.Text = "Visited Cites"
	g.X.Label.Text = "X"
    g.Y.Label.Text = "Y"

	lpLine, lpPoints, err := plotter.NewLinePoints(makeVisitedPoints(data))
	if err != nil {
		panic(err)
	}
	lpLine.Color = color.RGBA{B: 255, A: 255}
	//lpLine.LineStyle.Width = vg.Points(2)
	//lpPoints.Shape = draw.PyramidGlyph{}
	lpPoints.Color = color.RGBA{R: 80,G: 80,B: 80, A: 255}
	g.Add( lpLine, lpPoints)
	if err := g.Save(4*vg.Inch, 4*vg.Inch, fname); err != nil {
		panic(err)
	}
}
func SaveTourLengthImage(fname string,data1,data2 []TourLengthData){
	g, err := plot.New()
	if err != nil {
		panic(err)
	}

	g.Title.Text = "Tour length"
	g.X.Label.Text = "iteration"
    g.Y.Label.Text = "Length"
	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.
	//g.Y.Tick.Marker = commaTicks{}

	err = plotutil.AddLinePoints(g,
		"best", makePoints(data1),
		"avg", makePoints(data2))
	//	"second", randomPoints(len(data)))

	if err != nil {
		panic(err)
	}
	if err := g.Save(6*vg.Inch, 4*vg.Inch, fname); err != nil {
		panic(err)
	}
}
func makeVisitedPoints(data []CityData) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		//fmt.Println(i+1,data[i].X, data[i].Y)
		pts[i].X =data[i].X
		pts[i].Y =data[i].Y
	}
	return pts
}
func makePoints(data []TourLengthData) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X =float64(data[i].Iteration)
		pts[i].Y = float64(data[i].Data)
	}
	return pts
}
// RandomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = (pts[i].X + 10*rand.Float64()) * 1000
	}
	return pts
}

