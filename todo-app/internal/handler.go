package task

import (
	"html/template"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

func handleIndex(w http.ResponseWriter, r *http.Request) {
	TasksMutex.Lock()
	defer TasksMutex.Unlock()

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

	TasksMutex.Lock()
	tasks = append(tasks, Task{ID: nextID, Title: title})
	nextID++
	saveTasks()
	TasksMutex.Unlock()

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

	TasksMutex.Lock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Complete = true
			break
		}
	}
	saveTasks()
	TasksMutex.Unlock()

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

	TasksMutex.Lock()
	newTasks := []Task{}
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks
	saveTasks()
	TasksMutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
