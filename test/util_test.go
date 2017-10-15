package test

import (
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
	//t.pass()
  
}