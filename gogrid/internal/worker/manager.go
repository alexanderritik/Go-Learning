package worker

import (
	"fmt"

	"github.com/alexanderritik/gogrid/internal/store"
	"github.com/alexanderritik/gogrid/internal/task"
	"github.com/google/uuid"
)

type Manager struct {
	//Workeer Name -> Woker Struct
	Workers       map[string]*Worker
	TaskDB        store.Store          // Queue of task which is not assigned
	WorkerTaskMap map[uuid.UUID]string // TaskId -> WorkerName
}

func NewManager(workers []string, s store.Store) *Manager {

	workersMap := make(map[string]*Worker)

	for _, wName := range workers {
		workersMap[wName] = &Worker{
			Name:      wName,
			Queue:     make(chan task.Task, 10),
			Db:        make(map[uuid.UUID]*task.Task),
			TaskCount: 0,
			Store:     s,
		}
	}

	return &Manager{
		Workers:       workersMap,
		TaskDB:        s,
		WorkerTaskMap: make(map[uuid.UUID]string),
	}

}

func (m *Manager) GetPendingTasks() int {
	tasks, err := m.TaskDB.List()
	if err != nil {
		return 0
	}

	count := 0
	for _, t := range tasks {
		if t.State == task.Pending || t.State == task.Scheduled {
			count++
		}
	}
	return count
}
func (m *Manager) SelectWroker() string {
	var bestWorker string
	minTasks := 1000

	for k, v := range m.Workers {

		pendingTask := len(v.Queue)

		if pendingTask <= minTasks {
			bestWorker = k
			minTasks = pendingTask
		}
	}

	return bestWorker
}

func (m *Manager) AddTask(t task.Task) {

	m.TaskDB.Put(t)

	wName := m.SelectWroker()

	m.WorkerTaskMap[t.ID] = wName
	m.Workers[wName].Queue <- t

	fmt.Printf("Manager: Assigned task %s to worker %s\n", t.ID, wName)
}

func (m *Manager) UpdateTask(t task.Task) {
	m.TaskDB.Put(t)
}
