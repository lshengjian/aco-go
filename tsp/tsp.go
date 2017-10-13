
package tsp
import (
	"math"

)
type Matrix [][]float64

type City struct {
	X, Y float64
}

type TSP interface { //问题接口
	//Init()
	GetSize() int
	GetName() string
	GetLocations() []City
	GetDistanceMatrix() Matrix

}

//calculates edge weight (euclidiean distance)
func CalEdge(c1, c2 City) float64 {
	dx:=c2.X-c1.X
	dy:=c2.Y-c1.Y
	return math.Sqrt(dx*dx+dy*dy)
}
