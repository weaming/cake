package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron"
)

func In(a string, b []string) bool {
	for _, x := range b {
		if a == x {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) > 1 && In(os.Args[1], []string{"", "usage"}) {
		fmt.Println("See https://godoc.org/github.com/robfig/cron to learn how to define the cron field")
		os.Exit(0)
	}

	config := flag.String("config", "./cake.yml", "The config path in YAML format")
	detailLog := flag.Bool("log", false, "Print command stdout and stderr")
	flag.Parse()

	programs := ParseProgramConfig(*config)
	c := cron.New()
	c.Start()
	done := make(chan int)
	for i, p := range programs {
		log.Printf("The %v task is %+v\n", i, *p)
		err := c.AddFunc(p.Cron, func() {
			log.Printf("Running command %+v\n", p)
			exitCode, stdout, stderr := RunCommand(p.Envs, p.Dir, p.Command, p.Args...)
			if *detailLog || exitCode != 0 {
				text := fmt.Sprintf("Command result %v(%v %v):: exitCode: %v\n=== stdout ===\n%v\n===stderr===\n%v", p.Name, p.Command, p.Args, exitCode, stdout, stderr)
				log.Println(text)
			}

			WriteLog(p.StdoutFile, stdout)
			WriteLog(p.StderrFile, stderr)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}
