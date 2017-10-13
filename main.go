package main

import (
//	"log"
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
			Value: 10,
			Usage: "ant poplation.",
		},
		cli.IntFlag{
			Name:  "trys, t",
			Value: 100,
			Usage: "try times.",
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
		tries := c.Int("trys")
		fname := c.String("file_tsp")
		//log.Println("ants:",ants,"tsp-file:",fname)
		tsp:=tsp.NewFileTSP(fname)
		//int(arg_ants)
		swarm:=aco.NewColony(ants,tries,1,3,0.1,1,1,tsp)
		swarm.IsQuick=false
        swarm.Run()
		//start(arg_f)
		return nil
	}
	app.Run(os.Args)
/*	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Failed to run app with %s: %v", os.Args, err)
	}*/
}