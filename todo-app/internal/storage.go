package task

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	tasks      = []Task{}
	TasksMutex = sync.Mutex{}
	nextID     = 1
)

func LoadTasks() {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Println("Файл tasks.json не найден. Начинаем с пустого списка.")
		return
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		log.Println("Ошибка чтения tasks.json:", err)
	}
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
}

func saveTasks() {
	file, err := os.Create("tasks.json")
	if err != nil {
		log.Println("Не удалось сохранить задачи:", err)
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(tasks)
}
