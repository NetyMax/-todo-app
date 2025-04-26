package main

import (
	"log"
	"net/http"
	"todo-app/internal/task"
)

func main() {
	task.LoadTasks()

	http.HandleFunc("/", task.HandleIndex)
	http.HandleFunc("/add", task.HandleAddTask)
	http.HandleFunc("/complete", task.HandleCompleteTask)
	http.HandleFunc("/delete", task.HandleDeleteTask)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
