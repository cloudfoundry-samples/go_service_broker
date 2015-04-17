package web_server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	ac "github.com/xingzhou/go_service_broker/aws_client"
	"github.com/xingzhou/go_service_broker/module"
	"github.com/xingzhou/go_service_broker/utils"
)

const (
	DEFAULT_POLLING_INTERVAL_SECONDS = 10
)

type Controller struct {
	InstanceMap map[string]module.ServiceInstance
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
	vars := mux.Vars(r)

	serviceInstanceGuid := vars["service_instance_guid"]
	body, _ := ioutil.ReadAll(r.Body)

	var instance module.ServiceInstance
	err := json.Unmarshal(body, &instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	awsClient := ac.NewClient("us-east-1")
	vmId, err := awsClient.CreateInstance()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instance.InternalId = vmId
	instance.Id = serviceInstanceGuid

	lastOperation := module.LastOperation{
		State:                    "in progress",
		Description:              "creating service instance...",
		AsyncPollIntervalSeconds: DEFAULT_POLLING_INTERVAL_SECONDS,
	}

	instance.LastOperation = &lastOperation

	c.InstanceMap[instance.Id] = instance
	fmt.Println(instance)
	w.WriteHeader(http.StatusCreated)
	response := module.CreateServiceInstanceResponse{
		DashboardUrl:  "http://dashbaord_url",
		LastOperation: &lastOperation,
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))
}

func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
}

func (c *Controller) RemoveServiceInstance(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) Bind(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) UnBind(w http.ResponseWriter, r *http.Request) {

}
