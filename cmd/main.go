package main

import (
	"ToDoList/internal/controllers"
	database "ToDoList/internal/infrastructure"
	"ToDoList/internal/usecase"
	"ToDoList/server"
	"context"
	"fmt"
	"os"
)

func main() {

	os.Create("out/newfile.txt")

	ctx := context.Background()

	connstr := os.Getenv("CONN_STRING")
	conn, err := database.InitDatabase(ctx, connstr)
	if err != nil {
		panic(err)
	}
	repo := database.NewPostgresRepo(conn)
	fmt.Println("Connect Database Postgres")

	uc := usecase.NewUsecase(repo)
	defer conn.Close(ctx)

	controller := controllers.NewController(uc)
	server := server.NewHTTPServer(controller)
	
	go func(){
		if err := server.StartHttpServer(); err != nil {
		fmt.Println("Server starting error:", err)
	}	
	}()
	<-make(chan struct{})

}
