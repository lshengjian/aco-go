
package tsp
import (
	"math"

)
type Matrix [][]float64
type IntMatrix [][]int
type City struct {
	X, Y float64
}

type TSP interface { //问题接口
	//Init()
	GetSize() int
	GetName() string
	GetLocations() []City
	GetDistanceMatrix() IntMatrix

}

//calculates edge weight (euclidiean distance)
func CalEdge(c1, c2 City) int {
	dx:=c2.X-c1.X
	dy:=c2.Y-c1.Y
	return int(math.Sqrt(dx*dx+dy*dy)+0.5)
}
