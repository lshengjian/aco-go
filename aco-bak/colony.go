package aco-x
import (
	"math"

	"sync"
	"github.com/lshengjian/aco-go/util"
    "github.com/lshengjian/aco-go/tsp"
)
type Ants []*Ant
/*

func (as Ants)  Len() int{
	return len(as)
}
func (as Ants)  Less(i, j int) bool{
	return as[i].walkLength()<as[j].walkLength
}
func (as Ants)  Swap(i, j int){
	as[i],as[j]=as[j],as[i]
	as[i].idx,as[j].idx=as[j].idx,as[i].idx
}*/
type Message struct {
	ant *Ant
	walk []int
}
type Colony struct {
   Alpha float64
   Beta float64
   Pho float64
   Q float64
   BasePheromone float64
  // trail_0 float64
   pheromones tsp.Matrix
   totals tsp.Matrix
   Problem tsp.TSP
   size int
   popSize int
   maxIterations int
   Population Ants
   bestLength int
   best []int

   locks [][]sync.RWMutex
   IsQuick bool
   resultCh chan Message
   bestData []util.TourLengthData
   avgData []util.TourLengthData
}
func NewColony(popSize,maxIterations int,alpha,beta,pho,ip,q float64,p tsp.TSP)(*Colony) {
	rt:= &Colony{}
	rt.bestData = make([]util.TourLengthData, 0)
	rt.avgData = make([]util.TourLengthData, 0)
	rt.size=p.GetSize()
	rt.Alpha=alpha
	rt.Beta=beta
	rt.Pho=pho
	rt.Q=q
	rt.BasePheromone=ip
	rt.Problem=p
	rt.popSize=popSize
	rt.maxIterations=maxIterations
	rt.bestLength=999999999
	rt.Population=make(Ants,popSize)
	rt.pheromones=make(tsp.Matrix,rt.size)
	rt.totals=make(tsp.Matrix,rt.size)
	rt.init()
	return rt
}
func (p *Colony) init(){
	p.resultCh=make ( chan Message,p.popSize)
	//ip:=p.BasePheromone
	size:=p.size
	p.locks=make([][]sync.RWMutex,size)
	for i:=range p.locks{
		p.locks[i]=make([]sync.RWMutex,size)
		for j:=range p.locks[i] {
			var lk sync.RWMutex
			p.locks[i][j]=lk
		}
	}
	
	p.init_pheromone_trails(  )
	/*
	
	*/
	for  i:=0;i < p.popSize;i++ {
		ant:= NewAnt(p)
		ant.idx=i
		p.Population[i] =ant
  }
}
func (p *Colony) init_pheromone_trails(){
	ant:= NewAnt(p)
	ant.nn_tour() 
	p.BasePheromone = 1.0 /  float64(p.size*ant.walkLength ) 
	for i := range p.pheromones {
		p.pheromones[i] = make([]float64, p.size)
		p.totals[i] = make([]float64, p.size)
		for j := range p.pheromones[i] {
			p.pheromones[i][j] = p.BasePheromone //initialize to base pheromone = 1
			p.totals[i][j] = p.BasePheromone
		}
	}
}
func (p *Colony) GetBest() []int {
	return p.best
}

func (p *Colony) GetBestLength() int {
	return p.bestLength
}
func (p *Colony)	GetPheromon(i,j int) float64{
	if p.IsQuick {
		p.locks[i][j].RLock()
	    defer p.locks[i][j].RUnlock()
	}
	return p.pheromones[i][j] 
}
func (p *Colony)	GetTotal(i,j int) float64{
	if p.IsQuick {
		p.locks[i][j].RLock()
	    defer p.locks[i][j].RUnlock()
	}
	return p.totals[i][j] 
}
func (p *Colony)	SetPheromon(i,j int,data float64) {
	if p.IsQuick {
		p.locks[i][j].Lock()
	   defer p.locks[i][j].Unlock() 
	}
	//if data < 1e-10 {
	//	data=p.BasePheromone*1e-2
	//}
	p.pheromones[i][j] = data
	p.pheromones[j][i] = data
	data2:= math.Pow(data,p.Alpha) * math.Pow(p.HEURISTIC(i,j), p.Beta)
	p.totals[i][j]=data2
	p.totals[j][i]=data2

}

func (p *Colony)	Run() {
	
	if p.IsQuick {
		for  _,ant:=range p.Population {
			go	ant.Run()
		}
	}
	tourLength:=0
	for t:=1;t <= p.maxIterations;t++{
		total:=0
		for  _,ant:=range p.Population {
			if p.IsQuick {
				msg:=<-p.resultCh
				tourLength=p.findBest(t,msg)
				msg.ant.walkLength=tourLength
				p.layPheromones(tourLength,msg.walk)
			}else{
				ant.Tour(t)
				tourLength=p.findBest(t,Message{ant,ant.walk})
				ant.walkLength=tourLength
				p.layPheromones(tourLength,ant.walk)
			}
			total+=tourLength
		}
		p.evaporatePheromones()
		
		if(t%10==1){
			p.avgData=append(p.avgData,
				util.TourLengthData{t,total/p.popSize})
		}
		
	}
	//end := time.Now().UnixNano()
/*	ts := fmt.Sprintf("times:%.1f[ms]", float64(end-start)*1e-6)
	fmt.Println(ts)
	util.SaveTourLengthImage("tour-len.png",p.bestData,p.avgData)
	vd:=make ([]util.CityData,p.size+1)
	for i,item:=range vd{
		cidx:=p.best[i]
		item.Idx=cidx
		city:=p.Problem.GetLocations()[cidx]
		item.X,item.Y=city.X,city.Y
		vd[i]=item
	}

    util.SaveVisitedImage("visited.png",vd)*/
}
//global_acs_pheromone_update
func (p *Colony) layPheromones(walkLength int ,walk []int) {
	dq:= p.Q *3.0/ float64(walkLength)*float64(p.bestLength)/float64(walkLength)

	//
	for i := 1;i < len(walk);i++ { //起始结点出现两次
		k1,k2:=walk[i-1],walk[i]
		data:=p.GetPheromon(k1,k2)+dq
		p.SetPheromon(k1,k2,data)
  }
}
func (p *Colony) HEURISTIC(m,n int) float64{
	return (1.0 / (float64(p.Problem.GetDistanceMatrix()[m][n]) + 0.1))
}     

func (p *Colony) evaporatePheromones() {
	size:=p.size
    for x:=0; x < size;x++ {
      for y:=0; y < x;y++ {
		   dq:=(1 - p.Pho) * p.pheromones[x][y]
		   p.SetPheromon(x,y, dq)
		} 
	}
}

func (p *Colony) findBest(x int,msg Message) int{
    l:=p.CalculateWalkLength(msg.walk)
    if l < p.bestLength {
	//	str:=fmt.Sprintf("iterate:%d (ant:%d) -->%d",x,msg.ant.idx+1,l)
		p.bestLength=l
	//	fmt.Println(str)
		p.best=msg.walk
	    p.bestData=append(p.bestData,util.TourLengthData{x,p.bestLength})
   }
   return l
}

func (p *Colony) CalculateWalkLength(walk []int) int {
	distances := p.Problem.GetDistanceMatrix()
	sum := 0
	for i := 1; i < len(walk);i++ {
		sum += distances[walk[i-1]][walk[i]]
    }
    return sum
}
//lay Pheromones
/*
{	sort.Sort(p.Population)
	maxIdx:=int(float64(p.popSize)*0.9)
	for  _,ant:=range p.Population {
		dq:= p.Q / float64(ant.GetWalkLength())
		
		if p.IsQuick {
			ant.lock.RLock()
		}
		//fmt.Println("ant :",ant.idx+1)
		for i := 1;i < (p.size+1);i++ { //起始结点出现两次
			fmt.Println("walk city :",ant.idx+1)
			k1,k2:=ant.walk[i-1],ant.walk[i]
			old:=p.GetPheromon(k1,k2)
			if i<=maxIdx {
				p.SetPheromon(k1,k2,old+dq)
			}else{
				p.SetPheromon(k1,k2,old-dq)
			}
		}
		if p.IsQuick {
			ant.lock.RUnlock()
		}
	}
}*/
