package aco
import (
	"fmt"

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
   
   pheromones tsp.Matrix
   Problem tsp.TSP
   size int
   popSize int
   maxIterations int
   Population Ants
   bestLength int
   best []int

   locks [][]sync.RWMutex
   IsQuick bool
   IsMakeImage bool
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
	
	rt.init()
	return rt
}
func (p *Colony) init(){
	p.resultCh=make ( chan Message,p.popSize)
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
		ant.idx=i
		p.Population[i] =ant
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
func (p *Colony)	SetPheromon(i,j int,data float64) {
	if p.IsQuick {
		p.locks[i][j].Lock()
	   defer p.locks[i][j].Unlock() 
	}
	if data < 1e-10 {
		data=p.BasePheromone*1e-2
	}
	p.pheromones[i][j] = data
	p.pheromones[j][i] = data
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

/*  end := time.Now().UnixNano()
	ts := fmt.Sprintf("times:%.1f[ms]", float64(end-start)*1e-6)
	fmt.Println(ts)*/
	if p.IsMakeImage {
		fname:="./reports/"+p.Problem.GetName()
		util.SaveTourLengthImage(fname+"-best.png",p.bestData,p.avgData)
		vd:=make ([]util.CityData,p.size+1)
		for i,item:=range vd{
			cidx:=p.best[i]
			item.Idx=cidx
			city:=p.Problem.GetLocations()[cidx]
			item.X,item.Y=city.X,city.Y
			vd[i]=item
		}
	
		util.SaveVisitedImage(fname+"-visited.png",vd)
		fmt.Println("output images for:",fname)

	}

}

func (p *Colony) layPheromones(walkLength int ,walk []int) {
	dq:= p.Q / float64(walkLength)
	//if walkLength==p.bestLength {
		dq*=float64(p.bestLength)/float64(walkLength)*3
	//}
	for i := 1;i < len(walk);i++ { //起始结点出现两次
		k1,k2:=walk[i-1],walk[i]
		data:=p.GetPheromon(k1,k2)+dq
		p.SetPheromon(k1,k2,data)
  }
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

func (p *Colony) findBest(step int,msg Message) int{
    l:=p.CalculateWalkLength(msg.walk)
    if l < p.bestLength {
	//	str:=fmt.Sprintf("iterate:%d (ant:%d) -->%d",x,msg.ant.idx+1,l)
		p.bestLength=l
	//	fmt.Println(str)
		p.best=msg.walk
	    p.bestData=append(p.bestData,util.TourLengthData{step,p.bestLength})
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
