package web_server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/cloudfoundry-samples/go_service_broker/config"
	"github.com/cloudfoundry-samples/go_service_broker/model"
	"github.com/cloudfoundry-samples/go_service_broker/utils"
)

var (
	conf = config.GetConfig()
)

type Server struct {
	controller *Controller
}

func CreateServer() *Server {
	return &Server{controller: CreateController(loadServiceInstances(), loadServiceBindings())}
}

func (s *Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", s.controller.Catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.GetServiceInstance).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.CreateServiceInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.RemoveServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.Bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.UnBind).Methods("DELETE")

	http.Handle("/", router)

	fmt.Println("Server started, listening on port " + conf.Port + "...")
	http.ListenAndServe(":"+conf.Port, nil)
}

// private methods
func loadServiceInstances() map[string]*model.ServiceInstance {
	var serviceInstancesMap map[string]*model.ServiceInstance

	err := utils.ReadAndUnmarshal(&serviceInstancesMap, conf.DataPath, conf.ServiceInstancesFileName)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("WARNING: service instance data file '%s' does not exist: \n", conf.ServiceInstancesFileName)
		serviceInstancesMap = make(map[string]*model.ServiceInstance)
	}

	return serviceInstancesMap
}

func loadServiceBindings() map[string]*model.ServiceBinding {
	var bindingMap map[string]*model.ServiceBinding

	err := utils.ReadAndUnmarshal(&bindingMap, conf.DataPath, conf.ServiceBindingsFileName)
	if err != nil {
		fmt.Printf("WARNING: key map data file '%s' does not exist: \n", conf.ServiceBindingsFileName)
		bindingMap = make(map[string]*model.ServiceBinding)
	}

	return bindingMap
}
