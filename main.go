package main

import (
	"fmt"
	"runtime"
	"os"
	"github.com/urfave/cli"
	"github.com/lshengjian/aco-go/tsp"
	"github.com/lshengjian/aco-go/aco"
)
// Version holds the current app version
var Version = "1.0.0"

// go install -ldflags "-s -w"
func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
	   cli.StringFlag{
			Name:  "file_tsp ,f",
			Value: "./data/eil51.tsp",
			Usage: "TSP file `filename`",
		},
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
		//	Value: 1,
			Usage: "use multi CPU cores.",
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
	
		ants := c.Int("ants")
		tries := c.Int("tries")
		fname := c.String("file_tsp")
		//
		tsp:=tsp.NewFileTSP(fname)
		
		swarm:=aco.NewColony(ants,tries,1.0,5.0,0.1,1.0,1.0,tsp)
		//speed:=
		swarm.IsQuick=c.Bool("speed")
		fmt.Println("tries:",tries,"ants",ants)
		if (!swarm.IsQuick){
			//runtime.GOMAXPROCS(1)
		}else{
			fmt.Println("use CPU cores:",runtime.NumCPU())
		}
		  
		swarm.Run()
		fmt.Println("best:",swarm.GetBest())
		//start(arg_f)
		return nil
	}
	app.Run(os.Args)
/*	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Failed to run app with %s: %v", os.Args, err)
	}*/
}