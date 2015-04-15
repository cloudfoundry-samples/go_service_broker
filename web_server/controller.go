package web_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xingzhou/go_service_broker/module"
	"github.com/xingzhou/go_service_broker/utils"
)

type Controller struct {
}

func (c *Controller) Catalog(w http.ResponseWriter, r *http.Request) {
	templatePath := utils.GetPath([]string{"config", "catalog.json"})

	bytes, err := utils.ReadFile(templatePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var catalog module.Catalog

	err = json.Unmarshal(bytes, &catalog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(catalog)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) RemoveServiceInstance(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) Bind(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) UnBind(w http.ResponseWriter, r *http.Request) {

}
