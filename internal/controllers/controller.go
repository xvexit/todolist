package controllers

import (
	"ToDoList/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusBadRequest)
		return
	}

	if task.Name == "" || task.Text == "" {
		http.Error(w, "name and text are required", http.StatusBadRequest)
		return
	}

	if err := c.uc.AddTask(сtx, task.Name, task.Text); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)	// сделать чтобы вытаскивал из юзкейса ентити
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response", err)
		return
	}
}

func (c *Controller) HandleTaskList(w http.ResponseWriter, r *http.Request) {
	сtx, cancel := context.WithTimeout(r.Context(), time.Second*60)
	defer cancel()
	tasks, err := c.uc.TaskList(сtx)
	if err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(tasks, "", "    ")
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

func (c *Controller) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*60)
	defer cancel()

	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok || idstr == "" {
		http.Error(w, NewErrorDto("id parameter is missing").ToString(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, NewErrorDto("invalid id format").ToString(), http.StatusBadRequest)
		return
	}

	if err := c.uc.DelTask(ctx, id); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Deleted!")); err != nil {
		fmt.Println("failed to write HTTP response", err)
		return
	}
}

func (c *Controller) HandleDoTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*60)
	defer cancel()

	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok || idstr == "" {
		http.Error(w, NewErrorDto("id parameter is missing").ToString(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, NewErrorDto("invalid id format").ToString(), http.StatusBadRequest)
		return
	}

	if err := c.uc.DoTask(ctx, id); err != nil {
		http.Error(w, NewErrorDto(err.Error()).ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
