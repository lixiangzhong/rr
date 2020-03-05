package main

import (
	"log"

	"github.com/lixiangzhong/rr"
)

func main() {
	e, err := rr.NewEngine(rr.EngineConfig{
		BuildCmd: "go build -o _app",
		RunCmd:   "./_app",
		WatcherConfig: rr.WatcherConfig{
			ExcludeDir: []string{"node_modules", "dist", "vendor"},
			Ext:        []string{".go"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer e.Stop()
	err = e.Watch(".")
	if err != nil {
		log.Fatal(err)
	}
	e.Start()
	c := make(chan bool)
	<-c
}
