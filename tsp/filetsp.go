
package tsp
import (
	"sort"
	"math/rand"
	"bufio"
	"os"
	"strconv"
	"strings"
	"github.com/lshengjian/aco-go/util"
)
type NNData struct{
	idx int
	length int
}
type FileTSP struct {
	size        int
	name string
	locations []City
	distanceMatrix IntMatrix
	nnlist IntMatrix
}
func (t *FileTSP) GetSize() int{
	return t.size
}
func (t *FileTSP) GetName() string{
	return t.name
}
func (t *FileTSP) GetLocations() []City{
	return t.locations
}
func (t *FileTSP) GetDistanceMatrix() IntMatrix{
	return t.distanceMatrix
}
func NewFileTSP(fname string)(TSP) {
	var rt FileTSP
	rt.name=fname
	
	if strings.HasSuffix(fname,"tsp"){
		ss:=strings.Split(fname,".")
		rt.name=ss[0]
		//fmt.Println(rt.name)
		rt.initGraph("./data/"+fname)
	}else{
		ss:=strings.Split(fname,"-")
		n,_:=strconv.Atoi(ss[1])
		rt.randGraph(n)
	}
	return &rt
}
func (t *FileTSP) GetNNIdx(i,j int) int{
	return t.nnlist[i][j]
}
func (t *FileTSP) makeMatrix(){
	for i := range t.distanceMatrix {
	    nncities:=make([]NNData,t.size)
		t.distanceMatrix[i] = make([]int, t.size)
		t.nnlist[i] = make([]int, t.size)
		for j := range t.distanceMatrix[i] {
			t.distanceMatrix[i][j] = CalEdge(t.locations[i], t.locations[j])
		    nncities[j]=NNData{j,t.distanceMatrix[i][j]}
		}
		sort.Slice(nncities,func (a,b int) bool{
			return nncities[a].length<nncities[b].length
		})
		for j,d := range nncities {
			t.nnlist[i][j]=d.idx
		}

	}
}
//making graph
func (t *FileTSP) randGraph(n int) {
	t.locations = make([]City, n)
	t.size = len(t.locations)
	t.distanceMatrix = make(IntMatrix, t.size)
	t.nnlist=make(IntMatrix, t.size)
	for i := range t.locations {
		t.locations[i]= City{rand.Float64()*100.0, rand.Float64()*100.0}
	}
	t.makeMatrix()

	
}
//making graph
func (t *FileTSP) initGraph(fname string) {
	
	t.readFile(fname)
	t.size = len(t.locations)
	t.distanceMatrix = make(IntMatrix, t.size)
	t.nnlist=make(IntMatrix, t.size)
	t.makeMatrix()
}
func (t *FileTSP) readFile(fname string) {
	i, dim := 0, 0
	startFlag := false
	if file, err := os.Open(fname); err == nil {
		// make sure it gets closed
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			if strings.Contains(str, "DIMENSION") {
				dim = getDim(str)
				break
			}
		}
		t.locations = make([]City, dim)
		for scanner.Scan() {
			str := scanner.Text()
			if strings.Contains(str, "EOF") {
				break
			} else if startFlag {
				x, y := tokenize(str)
				if i < dim {
					t.locations[i] = City{x, y}
					i++
				} else {
					startFlag = false
				}
			} else if strings.Contains(str, "NODE_COORD_SECTION") {
				startFlag = true
			}
		}
		// check for errors
		if err = scanner.Err(); err != nil {
			util.PrintErrorAndExit(err)
		}
	} else {
		util.PrintErrorAndExit(err)
	}
}



//gets number of cities from the file
func getDim(str string) (dim int) {
	s := strings.Split(str, ":")
	num := strings.TrimLeft(s[1], " ")
	if dim, err := strconv.Atoi(num); err == nil {
		return dim
	} else {
		util.PrintErrorAndExit(err)
	}
	return 0
}
//tokenizes and converts to float
func tokenize(str string) (x, y float64) {
	s := strings.Split(str, " ")
	strX, strY := s[1], s[2]
	x, err := strconv.ParseFloat(strX, 64) //converts string to float64
	if err != nil {
		util.PrintErrorAndExit(err)
	}
	y, err = strconv.ParseFloat(strY, 64) //converts string to float64
	if err != nil {
		util.PrintErrorAndExit(err)
	}
	return x, y
}
