package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks []Task

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", GetTasks).Methods("GET")
	r.HandleFunc("/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	// Serve the HTML and JavaScript files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":9000", nil)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = fmt.Sprintf("%d", len(tasks)+1)
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, task := range tasks {
		if task.ID == params["id"] {
			tasks[i].Completed = !tasks[i].Completed
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, task := range tasks {
		if task.ID == params["id"] {
			tasks = append(tasks[:i], tasks[i+1:]...)
			json.NewEncoder(w).Encode("Task deleted")
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
