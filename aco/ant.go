package aco

import (

	"math"
//	"fmt"
	"math/rand"
)

type Ant struct {
	colony     *Colony
	walkLength float64
	base       int

	walk    []int
	visited []bool
	id int
}

func NewAnt(c *Colony) *Ant {
	rt := Ant{}
	rt.colony = c
	
	return &rt

}
func (p *Ant) SetWalk(data []int) {
  p.walk=make([]int,len(data))
  copy(p.walk,data)
}

/**
 * Set the base node index for this ant
 *
 * @param {Number} baseId
 */
func (p *Ant) setBase(base int) {
	p.base = base
	p.walkLength = 1e99
	p.walk = make([]int, 1)
	p.walk[0] = base
	//p.visited = nil
	p.visited = make([]bool,p.colony.Problem.GetSize())
	
/*	for i:=range p.visited{
		p.visited[i]=false
	}*/
	p.visited[base] = true
}


func (p *Ant) doWalk() {
	//size := p.colony.Problem.GetSize()
    p.setBase(rand.Intn(p.colony.size))
	for i := 1;i < p.colony.size;i++ {
		next := p.chooseNext(p.walk[i-1])
		p.walk = append(p.walk, next)
		p.visited[next] = true
	//	fmt.Println(p.id,"choose next:",next)
		
	}
	p.walk = append(p.walk, p.walk[0])
	p.walkLength = p.CalculateWalkLength()
	p.colony.wg.Done()
}

func (p *Ant) chooseNext(currentNode int) int {
	distances := p.colony.Problem.GetDistanceMatrix()
	pheromones := p.colony.Pheromones
	alpha := p.colony.Alpha
	beta := p.colony.Beta
	sum := 0.0
	unvisited := make([]int,0)
	//size := p.colony.Problem.GetSize()
	probs := make([]float64,0)
	for i:=range p.visited {
		//if _, ok := p.visited[i]; !ok && i != currentNode {
		if  !p.visited[i] && i != currentNode {
			unvisited = append(unvisited, i)
			data:= math.Pow(pheromones[currentNode][i], alpha) * math.Pow((1/distances[currentNode][i]), beta)
			sum += data
			probs=  append(probs,data)
		}
	}
	if len(unvisited)==1{
		return unvisited[0]
	}
   /* if len(unvisited)<1{
	   fmt.Println("chooseNext error")
	   return -1
    }*/
			
	for i:=range unvisited {
	   probs[i] =  probs[i]/sum
	}


	rnd := rand.Float64()
	x:=0
	tally := probs[0]
	for ;rnd > tally && x < len(unvisited)-1; x++{
		tally += probs[x]
	}
	return unvisited[x]
}

func (p *Ant) CalculateWalkLength() float64 {
	distances := p.colony.Problem.GetDistanceMatrix()
	sum := 0.0
	
	for i := 1; i < len(p.walk);i++ {
		sum += distances[p.walk[i-1]][p.walk[i]]
    }
    return sum
}


