package usecase

import (
	"ToDoList/internal/entity"
	"context"
)

type TaskListRepo interface {
	AddTask(ctx context.Context, t *entity.Task) error
	DeleteTask(ctx context.Context, id int64) error
	Update(ctx context.Context, t *entity.Task) error
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
	GetTaskById(ctx context.Context, id int) (*entity.Task, error)
	GetAllTasksPages(ctx context.Context, n int) (map[int][]entity.Task, error)
}
