package main

import (

	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/lshengjian/aco-go/aco"
	"github.com/lshengjian/aco-go/tsp"
	"github.com/lshengjian/aco-go/util"
	"github.com/urfave/cli"
)

// Version holds the current app version
var Version = "1.3.1"
func run(app *cli.App){
	app.Action = func(c *cli.Context) error {
		rand.Seed(time.Now().UnixNano())
		//fs:=[4]string{"eil51.tsp","R-20","R-50","R-100"}
		ants := c.Int("ants")
		tries := c.Int("tries")
		iterations := c.Int("iterations")
		flist, err := ioutil.ReadDir("./data")
		util.CheckError(err)

		fmt.Println("iterations:", iterations, "ants", ants)
		if c.Bool("speed") {
			fmt.Println("CPU cores:", runtime.NumCPU())
		}

		if tries < 1 {
			tries=1
		}
		total := len(flist)
	
		timers := make([][]float64, total)
		datas := make([][]float64, total)
		timeData := make([]*util.TestData, total)
		valueData := make([]*util.TestData, total)
	
		fnames:=make([]string,total)
		for i, f := range flist {
			timers[i] = make([]float64, tries)
			datas[i] = make([]float64, tries)
			p := tsp.NewFileTSP(f.Name())
			fnames[i]=p.GetName()
			for t := 0; t < tries; t++ {
				t1 := time.Now()
				swarm := aco.NewColony(ants, iterations, 1.0, 5.0, 0.1, 1.0, 1.0, p)
				swarm.IsQuick = c.Bool("speed")
				swarm.IsMakeImage=tries==1
				swarm.Run()
				timers[i][t] = time.Since(t1).Seconds()
				fmt.Println(swarm.Problem.GetName(),"-->", swarm.GetBestLength())
				datas[i][t] =float64( swarm.GetBestLength() )
			}
			timeData[i] = &util.TestData{fnames[i], "", timers[i]}
			valueData[i] = &util.TestData{fnames[i], "", datas[i]}
		

		}
		if tries > 1 {
			r1 := &util.ResultData{}
			r2 := &util.ResultData{}
			for _, d := range timeData {
				r1.Results = append(r1.Results, d)
			}
			for _, d := range valueData {
				r2.Results = append(r2.Results, d)
			}
			flag := ""
			if c.Bool("speed") {
				flag = "(Q)"
			}
			
			oname:=fmt.Sprintf("A%dT%d",ants,tries)+ flag + ".txt"
			r1.SaveDataToFile("./reports/T-" +oname )
			r2.SaveDataToFile("./reports/V-" + oname )
			fmt.Println("output report for:",oname)
		}
		return nil
	}
	app.Run(os.Args)
}
// GOOS=windows GOARCH=amd64 go install -ldflags "-s -w"
// GOOS=linux GOARCH=amd64 go install -ldflags "-s -w"
// GOOS=darwin GOARCH=amd64 go install -ldflags "-s -w"
func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "ants, a",
			Value: 20,
			Usage: "ant poplation.",
		},
		cli.IntFlag{
			Name:  "iterations, i",
			Value: 200,
			Usage: "iterations.",
		},
		cli.IntFlag{
			Name:  "tries, t",
			Value: 1,
			Usage: "try times.",
		},
		cli.BoolFlag{
			Name:  "speed, s",
			Usage: "use multi CPU cores.",
		},
	/*	cli.StringFlag{
			Name:  "output ,o",
			Value: "A20",
			Usage: "output file `filename`",
		},*/
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
    run(app)


}
