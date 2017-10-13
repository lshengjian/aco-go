package aco
import (
	"math/rand"
	"fmt"
	"time"
    "sync"
    "github.com/lshengjian/aco-go/tsp"
)
type Colony struct {
   Alpha float64
   Beta float64
   Pho float64
   Q float64
   BasePheromone float64
   
   pheromones tsp.Matrix
   Problem tsp.TSP
   size int
   popSize int
   maxIterations int
   Population []*Ant
   bestLength float64
	 wg sync.WaitGroup 

	locks [][]sync.RWMutex
	 
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
	rt.pheromones=make(tsp.Matrix,rt.size)
	
	rt.init()
	return rt
}
func (p *Colony) init(){
	
	ip:=p.BasePheromone
	size:=p.size
	p.locks=make([][]sync.RWMutex,size)
	for i:=range p.locks{
		p.locks[i]=make([]sync.RWMutex,size)
		for j:=range p.locks[i] {
			var lk sync.RWMutex
			p.locks[i][j]=lk

		}
	}
	for i := range p.pheromones {
		p.pheromones[i] = make([]float64, size)
		for j := range p.pheromones[i] {
			if i != j {
				p.pheromones[i][j] = ip //initialize to base pheromone = 1
			} else {
				p.pheromones[i][j] = 0.0
			}
		}
	}
	
	for  i:=0;i < p.popSize;i++ {
		ant:= NewAnt(p)
		ant.id=i+1
		
		p.Population[i] =ant
  }
}

func (p *Colony)	GetPheromon(i,j int) float64{
	if p.IsQuick {
		p.locks[i][j].RLock()
	    defer p.locks[i][j].RUnlock()
	}
	return p.pheromones[i][j] 
}
func (p *Colony)	SetPheromon(i,j int,data float64) {
	if p.IsQuick {
		p.locks[i][j].Lock()
	   defer p.locks[i][j].Unlock() 
	}
	p.pheromones[i][j] = data
	p.pheromones[j][i] = data
}
func (p *Colony)	Run() {
//	fmt.Println(p.Pho)
	rand.Seed(time.Now().UnixNano())  
	start := time.Now().UnixNano()
	if p.IsQuick {
		p.wg.Add(p.popSize)
		for  _,ant:=range p.Population {
			go	ant.Run()
		}
		p.wg.Wait()
	}else{
		for t:=1;t <= p.maxIterations;t++{
			for  _,ant:=range p.Population {
		       ant.Tour(t)
			}
		}
	}
	
	end := time.Now().UnixNano()
	ts := fmt.Sprintf("times:%.1f[ms]", float64(end-start)*1e-6)
	fmt.Println(ts)
	//fmt.Println(p.pheromones[0])
}

	

func (p *Colony)  Evaporatepheromones() {
	size:=p.size
    for x:=0; x < size;x++ {
      for y:=0; y < x;y++ {
		   dq:=(1 - p.Pho) * p.pheromones[x][y]
		   if dq<p.BasePheromone{
			 dq=p.BasePheromone
		   }
		   p.pheromones[x][y] = dq
		   p.pheromones[y][x] = p.pheromones[x][y]
		} 
  }
}
/*
func (p *Colony) findBest(x int) {
   for  _,ant:=range p.Population {
      if (ant.walkLength < p.bestLength) {
				p.bestLength = ant.walkLength
				msg:=fmt.Sprintf("iterate:%d ant(%d)-->%.1f",x,ant.id,p.bestLength)
				fmt.Println(msg)
			//	fmt.Println(ant.walk)
		//  strconv.Itoa(x) 
	    }
   }
}*/