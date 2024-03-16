package main

import (
	"github.com/patyukin/go-redis-streams/internal/cronjob"
	"log"
	"sync"
)

func main() {
	cj := cronjob.NewCronJob()
	var wg sync.WaitGroup

	err := cj.Add("* * * * *", func() {
		// TODO: set function
	})
	if err != nil {
		log.Fatalln("error with scheduler")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cj.Start()
	}()

	wg.Wait()
	cj.Stop()
}
