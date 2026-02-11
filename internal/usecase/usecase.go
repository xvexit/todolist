package usecase

import (
	"ToDoList/internal/entity"
	"context"
)

type Usecases struct {
	repo TaskListRepo
}

func NewUsecase(ctx context.Context, repo TaskListRepo) *Usecases {
	return &Usecases{
		repo: repo,
	}
}

func (uc *Usecases) AddTask(ctx context.Context, name, text string) error {
	task := entity.NewTask(name, text)
	if err := uc.repo.AddTask(ctx, task); err != nil {
		return err
	}
	return nil
}

func (uc *Usecases) DoTask(ctx context.Context, name string) error {

	task, err := uc.repo.GetTaskByName(ctx, name)
	if err != nil {
		return err
	}

	task.DoTask()

	if err := uc.repo.Update(ctx, task); err != nil {
		return err
	}

	return nil
}

func (uc *Usecases) DelTask(ctx context.Context, name string) error {

	task, err := uc.repo.GetTaskByName(ctx, name)
	if err != nil {
		return err
	}

	if err := uc.repo.DeleteTask(ctx, task); err != nil {
		return err
	}

	return nil
}

func (uc *Usecases) TaskList(ctx context.Context) ([]entity.Task, error) {
	return uc.repo.GetAllTasks(ctx)
}

func (uc *Usecases) GetAllTasksInPages(ctx context.Context, n int)(map[int][]entity.Task, error){
	return uc.repo.GetAllTasksPages(ctx, n)
}
