package test

import (
	"sort"
	"fmt"
	"math/rand"
    "math"
	"testing"


)

func Test00(t *testing.T) {
	sum:=0
	for i:=0 ;i<100;i++{
	   k:= rand.NormFloat64()
	   if math.Abs(k)<0.5{
		   sum+=1
	   }
	}
	fmt.Println("NormFloat64 :",sum)
/*	data:=[]NNData{NNData{0,3},NNData{1,2},NNData{2,4}}
	fmt.Println(data)
	sort.Slice(data,func (a,b int) bool{
		return data[a].length<data[b].length
	})
	fmt.Println(data)*/
  
}