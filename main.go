package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"log"
)


type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` 
}


var (
	thanoojTasks []Task
	idCounter    int
	mu           sync.Mutex
)


func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid task input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	idCounter++
	newTask.ID = idCounter
	newTask.Status = "pending" 
	thanoojTasks = append(thanoojTasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}


func GetTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thanoojTasks)
}


func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/tasks/"):] 
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, task := range thanoojTasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}


func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/tasks/"):] // Extract ID from the URL
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid task input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range thanoojTasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				task.Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				task.Description = updatedTask.Description
			}
			if updatedTask.Status == "pending" || updatedTask.Status == "completed" {
				task.Status = updatedTask.Status
			}

			thanoojTasks[i] = task

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}


func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/tasks/"):] // Extract ID from the URL
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range thanoojTasks {
		if task.ID == id {
			thanoojTasks = append(thanoojTasks[:i], thanoojTasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateTask(w, r)
		case http.MethodGet:
			GetTasks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetTaskByID(w, r)
		case http.MethodPut:
			UpdateTask(w, r)
		case http.MethodDelete:
			DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server starting on port 8081...")
if err := http.ListenAndServe(":8081", nil); err != nil {
    log.Fatalf("Failed to start server: %v", err)
}
}
