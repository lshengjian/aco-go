package test

import (
	//"fmt"

	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lshengjian/aco-go/tsp"
	"github.com/lshengjian/aco-go/aco"
)


func Test01ACO(t *testing.T) {
	a := assert.New(t)
	tsp:=tsp.NewFileTSP("../data/eil51.tsp")
	swarm:=aco.NewColony(1,1,1,4.5,0.16,1,1,tsp)
	
	a.Equal(51,tsp.GetSize())
	data:=[52]int{1,22,8,26,31,28,3,36,35,20,2,29,21,16,50,34,30,9,49,10,39,33,45,15,44,42,40,19,41,13,25,14,24,43,7,23,48,6,27,51,46,12,47,18,4,17,37,5,38,11,32,1}
	var idxs [52]int
	for i,d:=range data{
		idxs[i]=d-1
	}
	ant :=swarm.Population[0]
	ant.SetWalk(idxs[:])
	ll:=ant.CalculateWalkLength()
	//fmt.Println(ll)
	a.Equal(429,int(ll))
}




