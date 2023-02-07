package main

import (
	"fmt"
	"os"
	"time"

	"pt/command"
	"pt/config"
	"pt/store"

	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("app started")

	conf := config.Load()

	fmt.Printf("-- hosts: %q\n-- store file: %s\n-- debug: %t\n\n", conf.Hosts, conf.StoreFile, conf.Debug)

	sqlite, err := store.NewSqLite(conf.StoreFile)
	if err != nil {
		panic(err)
	}

	cmd := command.New(sqlite)

	ticker := time.NewTicker(5000 * time.Millisecond)

	for {
		select {
		case <-ticker.C:

			f := func() {
				if err := cmd.Run(conf.Hosts); err != nil {
					fmt.Printf("cmd run error: %v", err)
					os.Exit(1)
				}
			}

			// run only at 00 minutes of every hour
			if !conf.Debug && time.Now().Minute() == 0 {
				f()
			}

			if conf.Debug {
				f()
			}
		}
	}
}
