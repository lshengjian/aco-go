package test

import (


	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lshengjian/aco-go/tsp"
)

func Test01TSP(t *testing.T) {
	a := assert.New(t)
	c1:=tsp.City{0,0}
	c2:=tsp.City{3,4}
	a.Equal(5,tsp.CalEdge(c1,c2))
	tsp:=tsp.NewFileTSP("../data/eil51.tsp")
	a.Equal(51,tsp.GetSize())
	cities:=tsp.GetLocations()
	node:=cities[50]
	a.Equal(30.0,node.X)
	a.Equal(40.0,node.Y)
	dis:=tsp.GetDistanceMatrix()
//	fmt.Println(cities)
	a.Equal(12,dis[0][1])
}
func Test02TSP(t *testing.T) {
	a := assert.New(t)
	tsp:=tsp.NewFileTSP("../data/TSP_D")
	a.Equal(38,tsp.GetSize())
	cities:=tsp.GetLocations()
	node:=cities[0]
	 
	a.Equal(11003.611100,node.X)
	a.Equal(42102.500000,node.Y)

	//fmt.Println(cities)

}
func Test03TSP(t *testing.T) {
	a := assert.New(t)
	tsp:=tsp.NewFileTSP("../data/TSP_WS")
	a.Equal(29,tsp.GetSize())
	cities:=tsp.GetLocations()
	node:=cities[18]
	 
	a.Equal(26550.0000,node.X)
	a.Equal(13850.0000,node.Y)

	//fmt.Println(cities)

}





