package store

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/alexanderritik/gogrid/internal/task"
	"github.com/google/uuid"
)

type JsonStore struct {
	filename string
	mu       sync.RWMutex
}

func NewJsonStore(filename string) *JsonStore {
	return &JsonStore{filename: filename}
}

func (j *JsonStore) Get(id uuid.UUID) (task.Task, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	tasks, err := j.load()
	if err != nil {
		return task.Task{}, err
	}

	t, ok := tasks[id]
	if !ok {
		return task.Task{}, ErrTaskNotFound
	}
	return *t, nil
}

func (j *JsonStore) Put(t task.Task) error {
	j.mu.Lock()         // lock for writing
	defer j.mu.Unlock() // unlock when function finsihes

	tasks, err := j.load()

	if err != nil {
		return err
	}

	tasks[t.ID] = &t

	return j.writeBytes(tasks)
}

func (j *JsonStore) List() (map[uuid.UUID]*task.Task, error) {
	j.mu.RLock() // Read Lock (allows multiple readers, blocks writers)
	defer j.mu.RUnlock()
	return j.load()
}

func (j *JsonStore) load() (map[uuid.UUID]*task.Task, error) {
	f, err := os.OpenFile(j.filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	tasks := make(map[uuid.UUID]*task.Task)

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() == 0 {
		return tasks, nil // Return empty map if file is empty
	}

	decode := json.NewDecoder(f)
	if err := decode.Decode(&tasks); err != nil {
		return tasks, nil
	}

	return tasks, nil
}

func (j *JsonStore) writeBytes(tasks map[uuid.UUID]*task.Task) error {
	f, err := os.OpenFile(j.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	return encoder.Encode(tasks)
}
