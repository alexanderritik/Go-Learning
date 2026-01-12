package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexanderritik/gogrid/internal/task"
	"github.com/alexanderritik/gogrid/internal/worker"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Manager *worker.Manager
	Router  *chi.Mux
}

func NewApi(address string, port int, manager *worker.Manager) *Api {
	api := &Api{
		Address: address,
		Port:    port,
		Manager: manager,
		Router:  chi.NewRouter(),
	}

	api.initRoutes()
	return api
}

func (api *Api) initRoutes() {
	api.Router.Get("/task", api.GetTaskHandler)
	api.Router.Post("/task", api.StartTaskHandler)
	api.Router.Get("/stats", api.GetStatsHandler)
}

func (a *Api) Start() error {
	// Low level string formatting for address
	addr := fmt.Sprintf("%s:%d", a.Address, a.Port)
	log.Printf("API Server starting on %s", addr)
	return http.ListenAndServe(addr, a.Router)
}

func (api *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		PendingTasks int            `json:"pending_tasks"`
		RunningTasks int            `json:"running_tasks"` // Bonus metric
		Worker       map[string]int `json:"worker"`
	}

	result := Result{
		Worker: make(map[string]int),
	}

	result.PendingTasks = api.Manager.GetPendingTasks()
	for wName, w := range api.Manager.Workers {
		result.Worker[wName] = w.TaskCount
		result.RunningTasks += len(w.Queue)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}
func (api *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	t := task.New(body.Name, body.Image)

	api.Manager.AddTask(t)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (api *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := api.Manager.TaskDB.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
