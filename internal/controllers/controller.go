package controllers

import (
	"ToDoList/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Controller struct {
	uc *usecase.Usecases
}

func NewController(uc *usecase.Usecases) *Controller {
	return &Controller{
		uc: uc,
	}
}

func (c *Controller) HandlerAddTask(w http.ResponseWriter, r *http.Request) {
	сtx, cancel := context.WithTimeout(r.Context(), time.Second*60)
	defer cancel()
	var task TaskDto

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusBadRequest)
		return
	}

	if err := c.uc.AddTask(сtx, task.Name, task.Text); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response", err)
		return
	}
}

func (c *Controller) HandleTaskList(w http.ResponseWriter, r *http.Request){
	сtx, cancel := context.WithTimeout(r.Context(), time.Second*60)
	defer cancel()
	tasks, err := c.uc.TaskList(сtx)
	if err != nil{
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil{
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil{
		fmt.Println("failed to write HTTP response", err)
		return
	}
}