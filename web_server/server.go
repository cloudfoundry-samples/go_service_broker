package web_server

import (
	m "github.com/xingzhou/go_service_broker/module"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	controller *Controller
}

func CreateServer() *Server {
	return &Server{
		controller: &Controller{},
	}
}

func (s *Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", s.controller.Catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.CreateServiceInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.RemoveServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.Bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.Bind).Methods("DELETE")

	http.Handle("/", router)

	http.ListenAndServe(":8001", nil)
}

func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {

	catalog := m.Catalog{}

	w.Header().Set("Content-type", "application/json")
	data, _ := json.Marshal(catalog)
	fmt.Fprintf(w, string(data))
}
