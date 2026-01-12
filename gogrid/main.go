package main

import (
	"fmt"
	"os"

	"github.com/alexanderritik/gogrid/internal/api"
	"github.com/alexanderritik/gogrid/internal/store"
	"github.com/alexanderritik/gogrid/internal/worker"
)

func main() {
	os.Remove("tasks.json")

	fmt.Println("Starting GoGrid Orchestrator...")

	s := store.NewJsonStore("task.json")

	m := worker.NewManager([]string{"worker-1", "worker-2"}, s)

	for _, w := range m.Workers {
		go w.Run()
	}

	apiServer := api.NewApi("localhost", 5555, m)
	apiServer.Start()
}
