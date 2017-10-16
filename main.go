package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"os"
	"io/ioutil"
	"github.com/urfave/cli"
	"github.com/lshengjian/aco-go/util"
	"github.com/lshengjian/aco-go/tsp"
	"github.com/lshengjian/aco-go/aco"
)
// Version holds the current app version
var Version = "1.0.0"

// go install -ldflags "-s -w"
func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "ants, a",
			Value: 20,
			Usage: "ant poplation.",
		},
		cli.IntFlag{
			Name:  "tries, t",
			Value: 200,
			Usage: "try times.",
		},
		cli.BoolFlag{
			Name:  "speed, s",
			Usage: "use multi CPU cores.",
		},
		cli.StringFlag{
			Name:  "output ,o",
			Value: "./mydata.txt",
			Usage: "output file `filename`",
		},
	}
	app.Name = "ACO-GO"
	app.Authors = []cli.Author{
	    cli.Author{
			Name:  "Liu Shengjian",
			Email: "lsj178@139.com",
		},
	}
	app.Usage = "ACO demo"
	app.Version = Version
	
	app.Action = func(c *cli.Context) error {
		rand.Seed(time.Now().UnixNano())  
		//fs:=[4]string{"eil51.tsp","R-20","R-50","R-100"}
		
		ants := c.Int("ants")
		tries := c.Int("tries")
		flist,err:=ioutil.ReadDir("./data")
		util.CheckError(err)
		
		
		
		fmt.Println("tries:",tries,"ants",ants)
		if c.Bool("speed") {
			fmt.Println("CPU cores:",runtime.NumCPU())
		}
		
		cnt:=5
		timers:=make([][]float64,len(flist))
		data:=make([]*util.TestData,len(flist))
		for i,f:=range flist{
			timers[i]=make([]float64,cnt)
			p:=tsp.NewFileTSP(f.Name())
			for t := 0; t < cnt; t++ {
				t1 := time.Now()
				swarm:=aco.NewColony(ants,tries,1.0,5.0,0.1,1.0,1.0,p)
				swarm.IsQuick=c.Bool("speed")
				swarm.Run()
				timers[i][t]=time.Since(t1).Seconds()
				fmt.Println(t,"best:",swarm.GetBestLength())
				fmt.Println(t,"cost time:",timers[i][t])
			}
			data[i]=&util.TestData{p.GetName(),"",timers[i]}
		}
		
		r:=&util.ResultData{}
		for _,d :=range data	{
            r.Results=append(r.Results,d)
		}
		
		r.SaveDataToFile(c.String("output"))
		return nil
	}
	app.Run(os.Args)

}