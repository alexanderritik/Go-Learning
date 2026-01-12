package store

import (
	"errors"

	"github.com/alexanderritik/gogrid/internal/task"
	"github.com/google/uuid"
)

var (
	ErrTaskNotFound = errors.New("Task not found")
)

type Store interface {
	Put(t task.Task) error
	Get(id uuid.UUID) (task.Task, error)
	List() (map[uuid.UUID]*task.Task, error)
}
