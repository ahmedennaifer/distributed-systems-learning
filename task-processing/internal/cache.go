package internal

import (
	"fmt"
	"github.com/google/uuid"
)

type Cache struct {
	Db map[uuid.UUID]Task
}

func NewCache() *Cache {
	return &Cache{
		Db: make(map[uuid.UUID]Task, 0),
	}
}

func (c *Cache) Add(task Task) error {
	for _, ctask := range c.Db {
		if ctask.Name == task.Name {
			fmt.Printf("error: cannot add task %v, it already exists\n", task.Name)
			return fmt.Errorf("error: cannot add task %v, it already exists", task.Name)
		}
	}
	c.Db[task.ID] = task
	return nil
}

func (c *Cache) Get(taskID uuid.UUID) (Task, error) {
	task, exists := c.Db[taskID]
	if !exists {
		fmt.Printf("task %v does not exist\n", task.ID)
		return Task{}, fmt.Errorf("task %v does not exist\n", task.ID)
	}
	return task, nil
}
