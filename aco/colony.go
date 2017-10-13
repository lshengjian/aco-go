package aco
import (
	"math/rand"
	"fmt"
	"time"
 // "strconv"
  "github.com/lshengjian/aco-go/tsp"
)
type Colony struct {
   Alpha float64
   Beta float64
   Pho float64
   Q float64
   BasePheromone float64
   
   Pheromones tsp.Matrix
   Problem tsp.TSP
   size int
   popSize int
   maxIterations int
   Population []*Ant
	 bestLength float64
	 
	 IsQuick bool
}
func NewColony(popSize,maxIterations int,alpha,beta,pho,ip,q float64,p tsp.TSP)(*Colony) {
	rt:= &Colony{}
	rt.size=p.GetSize()
	rt.Alpha=alpha
	rt.Beta=beta
	rt.Pho=pho
	rt.Q=q
	rt.BasePheromone=ip
	rt.Problem=p
	rt.popSize=popSize
	rt.maxIterations=maxIterations
	rt.bestLength=1e99
	rt.Population=make([]*Ant,popSize)
	rt.Pheromones=make(tsp.Matrix,rt.size)

	rt.init()
	return rt
}
func (p *Colony) init(){
	ip:=p.BasePheromone
	size:=p.Problem.GetSize()
	for i := range p.Pheromones {
		p.Pheromones[i] = make([]float64, size)
		for j := range p.Pheromones[i] {
			if i != j {
				p.Pheromones[i][j] = ip //initialize to base pheromone = 1
			} else {
				p.Pheromones[i][j] = 0.0
			}
		}
	}
	
	for  i:=0;i < p.popSize;i++ {
		ant:= NewAnt(p)
		ant.id=i+1
		
		p.Population[i] =ant
  }
}
func (p *Colony) Step(idx int)  {

	//	fmt.Println("sendOutAnts")
	p.sendOutAnts()
	p.updatePheromones()
	p.findBest(idx)

}

func (p *Colony)	Run() {
//	fmt.Println(p.Pho)
	rand.Seed(time.Now().UnixNano())  
	start := time.Now().UnixNano()
	for i:=1;i <= p.maxIterations;i++{
		//fmt.Println(i)
		p.Step(i)
	}
	end := time.Now().UnixNano()
	ts := fmt.Sprintf("times:%.1f[ms]", float64(end-start)*1e-6)
	fmt.Println(ts)
	//fmt.Println(p.Pheromones[0])
}

func (p *Colony) sendOutAnts() {
	//size:=p.Problem.GetSize()
	for  _,ant:=range p.Population {
		ant.doWalk()
	}
}

func (p *Colony) updatePheromones() {
   p.evaporatePheromones()
   p.layPheromones()
}
	

func (p *Colony)  evaporatePheromones() {
	size:=p.Problem.GetSize()
    for x:=0; x < size;x++ {
      for y:=0; y < x;y++ {
		   dq:=(1 - p.Pho) * p.Pheromones[x][y]
		   if dq<p.BasePheromone{
			 dq=p.BasePheromone
		   }
		   p.Pheromones[x][y] = dq
		   p.Pheromones[y][x] = p.Pheromones[x][y]
		} 
  }
}
func (p *Colony) layPheromones() {
	pheromones := p.Pheromones
	for  _,ant:=range p.Population {
		dq:=(1 / ant.walkLength) * p.Q
	  for i := 1;i < len(ant.walk);i++ { //起始结点出现两次
		 pheromones[ant.walk[i-1]][ant.walk[i]] += dq
		  pheromones[ant.walk[i]][ant.walk[i-1]] +=dq
		}
  }
}
func (p *Colony) findBest(x int) {
   for  _,ant:=range p.Population {
      if (ant.walkLength < p.bestLength) {
				p.bestLength = ant.walkLength
				msg:=fmt.Sprintf("iterate:%d ant(%d)-->%.1f",x,ant.id,p.bestLength)
				fmt.Println(msg)
				fmt.Println(ant.walk)
		//  strconv.Itoa(x) 
	    }
   }
}