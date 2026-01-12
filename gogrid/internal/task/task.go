package task

import (
	"time"

	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID           uuid.UUID
	Name         string
	State        State
	Image        string
	Memory       int
	Disk         int
	ExposedPorts map[string]string
	StartTime    time.Time
	FinishTime   time.Time
}

func New(name, image string) Task {
	return Task{
		ID:    uuid.New(),
		Name:  name,
		State: Pending,
		Image: image,
	}
}
