package web_server

import (
	m "github.com/xingzhou/go_service_broker/module"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
}

func (s *Server) Start() {
	http.HandleFunc("/", s.catalog)
	http.ListenAndServe(":8001", nil)
}

func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	mux.NewRouter()

	catalog := m.Catalog{}

	w.Header().Set("Content-type", "application/json")
	data, _ := json.Marshal(catalog)
	fmt.Fprintf(w, string(data))
}
