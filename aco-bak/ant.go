package aco-x

import (
//	"math"
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

func (p *Ant) local_acs_pheromone_update(phase int)  {
	h, j:=p.walk[phase-1],p.walk[phase]
//	pher:=p.colony.GetPheromon(h,j)
// math.Pow(pher,p.colony.Alpha) * math.Pow(p.colony.HEURISTIC(h,j), p.colony.Beta)
	data:=p.colony.GetTotal(h,j)
	

   /* still additional parameter has to be introduced */
   p.colony.SetPheromon(h,j,(1. - 0.1) * data + 0.1 * p.colony.BasePheromone)

}
func (p *Ant) setBase(base int) {
	p.base = base
	p.walk = make([]int, 1)
	p.walk[0] = base
	p.visited = make([]bool,p.colony.size)
	p.visited[base] = true
}
func (p *Ant) nn_tour() {
	p.setBase(rand.Intn(p.colony.size))
	for i := 1;i < p.colony.size;i++ {
  	  next := p.choose_closest_next(p.walk[i-1])
	  p.walk = append(p.walk, next)
	  p.visited[next] = true
	}
	p.walk = append(p.walk, p.walk[0])
	p.walkLength=p.colony.CalculateWalkLength(p.walk)
}
func (p *Ant) Tour(t int)  {
	p.setBase(rand.Intn(p.colony.size))
	for i := 1;i < p.colony.size;i++ {
  	  next := p.chooseNext(p.walk[i-1])
	  p.walk = append(p.walk, next)
	  p.visited[next] = true
	 // p.local_acs_pheromone_update(i)
	}
	p.walk = append(p.walk, p.walk[0])
	//p.local_acs_pheromone_update(p.colony.size)
}
    

func (p *Ant) Run() {
	for t:=1;t <= p.colony.maxIterations;t++{
		p.Tour(t)
	    data:=make([]int,p.colony.size+1)
		copy(data,p.walk)
		p.colony.resultCh<-Message{p,data}
	} 
}
func (p *Ant) choose_closest_next(currentNode int) int {
	//distances := p.colony.Problem.GetDistanceMatrix()
	
	min_distance:=9999999
	min_idx:=0
	dis:=p.colony.Problem.GetDistanceMatrix()
	for i:=range p.visited {
		if  !p.visited[i] && i != currentNode {
			if dis[currentNode][i]<min_distance{
				min_distance=dis[currentNode][i]
				min_idx=i
			}
		}
	}
	return min_idx
}
func (p *Ant) chooseNext(currentNode int) int {
	//distances := p.colony.Problem.GetDistanceMatrix()
	//alpha := p.colony.Alpha
	//beta := p.colony.Beta
	sum := 0.0
	unvisited := make([]int,0)
	probs := make([]float64,0)
	best:=0.0
	bestIdx:=0
	for i:=range p.visited {
		if  !p.visited[i] && i != currentNode {
			//pher:=p.colony.GetPheromon(currentNode,i)
			unvisited = append(unvisited, i)
			data:= p.colony.GetTotal(currentNode,i)
			// math.Pow(pher,alpha) * math.Pow(1.0/float64(distances[currentNode][i]), beta)
			sum += data
			if data>best{
				best=data
				bestIdx=i
			}
			probs = append(probs,data)
		}
	}
	
	if len(unvisited)==1{
		return unvisited[0]
	}
	rnd := rand.Float64()
	if rnd<0.05{
		return bestIdx
	}
	
	ps:=0.0
	for i:=range unvisited {
		ps = ps+ probs[i]/sum
		probs[i]=ps
	}
	

	rt:=0
	for x:=0;x<len(probs); x++{
		if rnd<=probs[x]{
			rt=x
			break
		}
	}

	return unvisited[rt]
}



