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
   nn_ls int
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
	rt.nn_ls=rt.size/2
	if rt.nn_ls>24{
		rt.nn_ls=24
	}
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

/*    
 
      FUNCTION:       2-opt a tour 
      INPUT:          pointer to the tour that undergoes local optimization
      OUTPUT:         none
      (SIDE)EFFECTS:  tour is 2-opt
      COMMENTS:       the neighbourhood is scanned in random order (this need 
                      not be the best possible choice). Concerning the speed-ups used 
		      here consult, for example, Chapter 8 of
		      Holger H. Hoos and Thomas Stuetzle, 
		      Stochastic Local Search---Foundations and Applications, 
		      Morgan Kaufmann Publishers, 2004.
		      or some of the papers online available from David S. Johnson.
*/
func (p *Colony) two_opt_first(ant *Ant ) {
  var  c1, c2 int             /* cities considered for an exchange */
  var  s_c1, s_c2 int        /* successor cities of c1 and c2     */
  var  p_c1, p_c2 int         /* predecessor cities of c1 and c2   */   
  var  pos_c1, pos_c2 int     /* positions of cities c1, c2        */
  var  i, j, h, l ,radius int
  var  pos []int               /* positions of cities in tour */ 
  var  dlb []bool               /* vector containing don't look bits */ 
  n:=p.size
  pos = make([]int,n)
  dlb = make([]bool,n)
  distances:=p.Problem.GetDistanceMatrix()
  help, n_improves, n_exchanges:=0,0,0
  h1, h2, h3, h4,gain:=0,0,0,0,0

  for  i = 0 ; i < n ; i++  {
    pos[ant.walk[i]] = i
	  dlb[i] = false
  }

  improvement_flag := true
  random_vector := util.GenerateRandomPermutation( n )
  //fmt.Println("nn_ls:",p.nn_ls)
//  for  improvement_flag {
	for  k:=0;k<n*10;k++ {
    improvement_flag = false
 	for l = 0 ; l < n; l++ {
	  c1 = random_vector[l] 
	 // fmt.Println(l,c1)
	  if   dlb[c1] {
	    continue
	  }
	  pos_c1 = pos[c1]
	  s_c1 = ant.walk[pos_c1+1]
	  radius = distances[c1][s_c1]

	    /* First search for c1's nearest neighbours, use successor of c1 */
	  for  h = 0 ; h < p.nn_ls ; h++  {
	    c2 =  p.Problem.GetNNIdx(c1,h) /* exchange partner, determine its position */
	   // fmt.Println(c1,h,c2) 
	    if  radius > distances[c1][c2] {
	  	  s_c2 = ant.walk[pos[c2]+1]
		  gain =  - radius + distances[c1][c2] + 
					distances[s_c1][s_c2] - distances[c2][s_c2]
		  if ( gain < 0 ) {
			h1,h2,h3, h4  = c1,s_c1,c2,s_c2
			goto exchange2opt
		  }
		} else {
		  break
		}
	  }
	          
	    /* Search one for next c1's h-nearest neighbours, use predecessor c1 */
	  if pos_c1 > 0 {
	  	p_c1 = ant.walk[pos_c1-1]
	  }else{
	  	p_c1 =  ant.walk[n-1]
	  } 
	  radius =distances[p_c1][c1]
	  for  h = 0 ; h < p.nn_ls ; h++  {
		c2 = p.Problem.GetNNIdx(c1,h)  /* exchange partner, determine its position */
		if  radius > distances[c1][c2] {
		  pos_c2 = pos[c2]
		  if (pos_c2 > 0){
			p_c2 = ant.walk[pos_c2-1]
		  }else{
			p_c2 = ant.walk[n-1]
		  } 
		  if   p_c2 == c1|| p_c1 == c2 {
			continue
		  }
		  gain =  - radius + distances[c1][c2] + 
					distances[p_c1][p_c2] - distances[p_c2][c2]
		  if  gain < 0  {
			h1,h2,h3,h4 = p_c1,c1 ,p_c2 ,c2
			goto exchange2opt
		  }
		}else{
		  break
		}
	  }      
				/* No exchange */
      dlb[c1] = true
      continue
 exchange2opt:
	  n_exchanges++
	  improvement_flag = true
	  dlb[h1], dlb[h2] = false , false
	  dlb[h3] ,dlb[h4] = false, false
	
	  if  pos[h3] < pos[h1]  {
		h1, h3 = h3,h1
		h2, h4= h4,h2
	  }
	  if  pos[h3] - pos[h2] < n / 2 + 1 {
				/* reverse inner part from pos[h2] to pos[h3] */
		i ,j= pos[h2],pos[h3]
		for i < j {
		  c1 = ant.walk[i]
		  c2 = ant.walk[j]
		  ant.walk[i] = c2
		  ant.walk[j] = c1
		  pos[c1] = j
		  pos[c2] = i
		  i++
		  j--
		}
	  }else {
				/* reverse outer part from pos[h4] to pos[h1] */
		i ,j= pos[h1], pos[h4]
		if ( j > i ){
		  help = n - (j - i) + 1
		}else{
		  help = (i - j) + 1
		} 
		help = help / 2
		for  h = 0 ; h < help ; h++  {
		  c1 = ant.walk[i]
		  c2 = ant.walk[j]
		  ant.walk[i] = c2
		  ant.walk[j] = c1
		  pos[c1] = j
		  pos[c2] = i
		  i--
		  j++
		  if  i < 0 {
			i = n-1
		  }
		  if   j >= n {
			j = 0
		  }
		}
		ant.walk[n] = ant.walk[0]
	  }
	}
	if  improvement_flag  {
		n_improves++
	}
  }
}
