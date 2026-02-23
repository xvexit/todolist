package server

import (
	"ToDoList/internal/controllers"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	handlers *controllers.Controller
}

func NewHTTPServer(handlers *controllers.Controller) *HTTPServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

func (s *HTTPServer) StartHttpServer() error {

	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.handlers.HandlerAddTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.handlers.HandleTaskList)

	fmt.Println("Server starting on :5050...")
	err := http.ListenAndServe(":5050", router)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
