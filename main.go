package main

import (
	"log"

	"github.com/Seann-Moser/hypr-u/filemonitor"
)

func main() {
	cfg, err := filemonitor.LoadDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := filemonitor.NewFileWatcher(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// start the auto-restart worker
	watcher.StartWorker()

	log.Println("Watching for file changes...")
	watcher.Start()
}
