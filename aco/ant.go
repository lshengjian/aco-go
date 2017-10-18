package aco

import (
	"math"
	"math/rand"
)

type Ant struct {
	colony     *Colony
	base       int
    walkLength int //only update by colony!
	walk    []int
	visited []bool
	idx int
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


func (p *Ant) setBase(base int) {
	p.base = base
	p.walk = make([]int, 1)
	p.walk[0] = base
	p.visited = make([]bool,p.colony.size)
	p.visited[base] = true
}

func (p *Ant) Tour(t int) {
	p.setBase(rand.Intn(p.colony.size))
	for i := 1;i < p.colony.size;i++ {
  	  next := p.chooseNext(p.walk[i-1])
	  p.walk = append(p.walk, next)
	  p.visited[next] = true
	}
	p.walk = append(p.walk, p.walk[0])
    //p.walkLength=p.calculateWalkLength()
}

func (p *Ant) Run() {
	for t:=1;t <= p.colony.maxIterations;t++{
		p.Tour(t)
	    data:=make([]int,p.colony.size+1)
		copy(data,p.walk)
		p.colony.resultCh<-Message{p,data}
	} 
}

func (p *Ant) chooseNext(currentNode int) int {
	distances := p.colony.Problem.GetDistanceMatrix()
	alpha := p.colony.Alpha
	beta := p.colony.Beta
	sum := 0.0
	unvisited := make([]int,0)
	probs := make([]float64,0)
	for i:=range p.visited {
		if  !p.visited[i] && i != currentNode {
			pher:=p.colony.GetPheromon(currentNode,i)
			unvisited = append(unvisited, i)
			data:= math.Pow(pher,alpha) * math.Pow(1.0/float64(distances[currentNode][i]), beta)
			sum += data
			probs = append(probs,data)
		}
	}
	if len(unvisited)==1{
		return unvisited[0]
	}
	ps:=0.0
	for i:=range unvisited {
		ps = ps+ probs[i]/sum
		probs[i]=ps
	}
	rnd := rand.Float64()

	rt:=0
	for x:=0;x<len(probs); x++{
		if rnd<=probs[x]{
			rt=x
			break
		}
	}

	return unvisited[rt]
}


