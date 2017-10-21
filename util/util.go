package util

//GLPATH=/usr/lib make
import (
	"image/color"
	"os"
	"fmt"
	"math/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io/ioutil"
	"math"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	
)
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
	
/*  
	 Copyright (C) 2004  Thomas Stuetzle  
      FUNCTION:       generate a random permutation of the integers 0 .. n-1
      INPUT:          length of the array
      OUTPUT:         pointer to the random permutation
      (SIDE)EFFECTS:  the array holding the random permutation is allocated in this 
                      function. Don't forget to free again the memory!
      COMMENTS:       only needed by the local search procedures
*/
func GenerateRandomPermutation( n int)([]int){
		i, help, node, tot_assigned := 0,0,0,0
		r :=make ([]int,n)
	
		for  i = 0 ; i < n; i++{
		  r[i] = i
		} 
	
	   for  i = 0 ; i < n ; i++  {
		 /* find (randomly) an index for a free unit */ 
		 rnd  :=rand.Float64()
		 node = int (rnd  * float64((n - tot_assigned)) )
		 help = r[i]
		 r[i] = r[i+node]
		 r[i+node] = help
		 tot_assigned++
	   }
	   return r
	}

	
	func WriteFile(fname ,data string)  {
		err:=ioutil.WriteFile(fname,[]byte(data),0644)
		CheckError(err) 
	}
	type TestData struct{
		Name string
		Args string
		Data []float64
	}
	
	func (p *TestData) Min() float64{
		return floats.Min(p.Data)
	}
	func (p *TestData) Max() float64{
		return floats.Max(p.Data)
	}
	func (p *TestData) Mean() float64{
		return stat.Mean(p.Data, nil)
	}
	func (p *TestData) Std() float64{
		return math.Sqrt(stat.Variance(p.Data, nil))
	}
	
	type ResultData struct{
	//	ProblemName string
	 //   Tries int
		Results []*TestData
	}
	func (p *ResultData) SaveDataToFile(fname string)  {
		 //fname:=p.ProblemName+".txt"
		 str:="TSP\tMean\tMin\tMax\tStd\n"
		 for _,d:=range p.Results{
			 str+=fmt.Sprintf("%s\t%.4G\t%.4G\t%.4G\t%.4G\n",d.Name,d.Mean(),d.Min(),d.Max(),d.Std())
		 }
		 WriteFile(fname,str)
	
	}
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

