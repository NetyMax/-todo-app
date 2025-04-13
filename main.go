package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Task представляет одну задачу
type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Complete bool   `json:"complete"`
}

// Хранилище задач (в памяти)
var (
	tasks      = []Task{}
	tasksMutex = sync.Mutex{} // Защита от гонки при доступе из нескольких горутин
	nextID     = 1
)

// Загрузка задач из файла
func loadTasks() {
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
	// Обновляем nextID
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
}

// Сохранение задач в файл
func saveTasks() {
	file, err := os.Create("tasks.json")
	if err != nil {
		log.Println("Не удалось сохранить задачи:", err)
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(tasks)
}

// Главная страница (HTML)
func handleIndex(w http.ResponseWriter, r *http.Request) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	tmpl.Execute(w, tasks)
}

// Добавление задачи
func handleAddTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Пустое название задачи", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	tasks = append(tasks, Task{ID: nextID, Title: title})
	nextID++
	saveTasks()
	tasksMutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Отметить задачу выполненной
func handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Complete = true
			break
		}
	}
	saveTasks()
	tasksMutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Удаление задачи
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	newTasks := []Task{}
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks
	saveTasks()
	tasksMutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	loadTasks()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAddTask)
	http.HandleFunc("/complete", handleCompleteTask)
	http.HandleFunc("/delete", handleDeleteTask)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HTML-шаблон (встроен прямо в код)
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Задачи</title>
</head>
<body>
    <h1>Мои задачи</h1>
    <form action="/add" method="POST">
        <input type="text" name="title" placeholder="Новая задача">
        <button type="submit">Добавить</button>
    </form>
    <ul>
        {{range .}}
        <li>
            {{if .Complete}}✅ {{end}}
            {{.Title}}
            {{if not .Complete}} <a href="/complete?id={{.ID}}">[выполнить]</a>{{end}}
            <a href="/delete?id={{.ID}}">[удалить]</a>
        </li>
        {{else}}
        <li>Нет задач</li>
        {{end}}
    </ul>
</body>
</html>
`
