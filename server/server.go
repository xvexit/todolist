package server

import (
	"ToDoList/internal/controllers"
	"errors"
	"fmt"
	"log"
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

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("REQUEST → %s %s from %s",
				r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(s.handlers.HandleDoTask)
	router.Path("/tasks").Methods("POST").HandlerFunc(s.handlers.HandlerAddTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.handlers.HandleTaskList)
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(s.handlers.HandleDeleteTask)

	fmt.Println("Server starting on :5050...")
	err := http.ListenAndServe(":5050", router)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
