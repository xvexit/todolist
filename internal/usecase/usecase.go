package usecase

import (
	"ToDoList/internal/entity"
	"context"
)

type Usecases struct {
	repo TaskListRepo
}

func NewUsecase(repo TaskListRepo) *Usecases {
	return &Usecases{
		repo: repo,
	}
}

func (uc *Usecases) AddTask(ctx context.Context, name, text string) error {
	task := entity.NewTask(name, text, 0)
	if err := uc.repo.AddTask(ctx, task); err != nil {
		return err
	}
	return nil
}

func (uc *Usecases) DoTask(ctx context.Context, id int64) error {

	task, err := uc.repo.GetTaskById(ctx, id)
	if err != nil {
		return err
	}

	task.DoTask()

	if err := uc.repo.Update(ctx, task); err != nil {
		return err
	}

	return nil
}

func (uc *Usecases) DelTask(ctx context.Context, id int64) error {
	return uc.repo.DeleteTask(ctx, id)
}

func (uc *Usecases) TaskList(ctx context.Context) ([]entity.Task, error) {
	return uc.repo.GetAllTasks(ctx)
}

func (uc *Usecases) GetAllTasksInPages(ctx context.Context, n int) (map[int][]entity.Task, error) {
	return uc.repo.GetAllTasksPages(ctx, n)
}

