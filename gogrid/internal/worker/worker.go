package worker

import (
	"fmt"
	"time"

	"github.com/alexanderritik/gogrid/internal/store"
	"github.com/alexanderritik/gogrid/internal/task"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     chan task.Task
	Db        map[uuid.UUID]*task.Task
	TaskCount int
	Store     store.Store
}

func (w *Worker) Run() {
	fmt.Printf("Worker %s: Starting\n", w.Name)

	for t := range w.Queue {
		fmt.Printf("Woker %s Found task %s\n", w.Name, t.ID)
		w.RunTask(t)
		fmt.Printf("Worker %s Finished task %s\n", w.Name, t.ID)
	}

}

func (w *Worker) RunTask(t task.Task) {

	t.State = task.Running
	t.StartTime = time.Now()
	w.Db[t.ID] = &t

	w.Store.Put(t)

	time.Sleep(5 * time.Second)

	t.FinishTime = time.Now()
	t.State = task.Completed
	w.Db[t.ID] = &t
	w.Store.Put(t)

	w.TaskCount++
}
