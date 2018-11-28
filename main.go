package main

import "fmt"

// import (
// 	"fmt"

// 	"github.com/robfig/cron"
// )

// func main() {
// 	c := cron.New()
// 	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
// 	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
// 	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
// 	c.Stop() // Stop the scheduler (does not stop any jobs already running).
// }

func main() {
	programs := ParseProgramConfig("./example.yaml")
	for _, p := range programs {
		fmt.Printf("%+v\n", *p)
	}
}
