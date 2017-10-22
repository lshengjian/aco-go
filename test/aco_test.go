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
	tsp:=tsp.NewFileTSP("eil51.tsp")
	swarm:=aco.NewColony(1,1,1,4.5,0.16,1,1,tsp)
	
	a.Equal(51,tsp.GetSize())
//	data:=[52]int{1,22,8,26,31,28,3,36,35,20,2,29,21,16,50,34,30,9,49,10,39,33,45,15,44,42,40,19,41,13,25,14,24,43,7,23,48,6,27,51,46,12,47,18,4,17,37,5,38,11,32,1}
    data:=[52]int{1 ,32 ,11 ,38 ,5 ,49 ,9 ,50 ,34 ,30 ,10 ,39 ,33 ,45 ,15 ,44 ,37 ,17 ,4 ,18 ,47 ,12 ,46 ,51 ,27 ,6 ,14 ,25 ,13 ,41 ,19 ,42 ,40 ,24 ,43 ,7 ,23 ,48 ,8 ,26 ,31 ,28 ,3 ,20 ,35 ,36 ,29 ,21 ,16 ,2 ,22 ,1}	
    var idxs [52]int
	for i,d:=range data{
		idxs[i]=d-1
	}
	ll:=swarm.CalculateWalkLength(idxs[:])
	//a.Equal(426,int(ll))
	//fmt.Println(ll)
	a.Equal(461,int(ll))
}




