package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func resetState() {
	tasks = []Task{}
	nextID = 1
	_ = os.Remove("tasks.json")
}

func TestHandleAddTask(t *testing.T) {
	resetState()

	form := url.Values{}
	form.Add("title", "Test Task")

	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handleAddTask(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect status, got %d", w.Code)
	}

	if len(tasks) != 1 || tasks[0].Title != "Test Task" {
		t.Errorf("Task was not added correctly: %+v", tasks)
	}
}

func TestHandleCompleteTask(t *testing.T) {
	resetState()
	tasks = append(tasks, Task{ID: 1, Title: "Incomplete Task"})

	req := httptest.NewRequest("GET", "/complete?id=1", nil)
	w := httptest.NewRecorder()

	handleCompleteTask(w, req)

	if !tasks[0].Complete {
		t.Errorf("Expected task to be marked as complete")
	}
}

func TestHandleDeleteTask(t *testing.T) {
	resetState()
	tasks = append(tasks, Task{ID: 1, Title: "To Be Deleted"})

	req := httptest.NewRequest("GET", "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handleDeleteTask(w, req)

	if len(tasks) != 0 {
		t.Errorf("Expected task to be deleted")
	}
}

func TestHandleIndex(t *testing.T) {
	resetState()
	tasks = append(tasks, Task{ID: 1, Title: "Test View Task"})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Test View Task")) {
		t.Errorf("Task title not found in output")
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	resetState()
	tasks = append(tasks, Task{ID: 1, Title: "Persisted Task"})
	saveTasks()
	tasks = nil
	loadTasks()

	if len(tasks) != 1 || tasks[0].Title != "Persisted Task" {
		t.Errorf("Task was not persisted properly")
	}
}
