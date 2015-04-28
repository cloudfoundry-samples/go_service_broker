package web_server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/xingzhou/go_service_broker/config"
	"github.com/xingzhou/go_service_broker/module"
	"github.com/xingzhou/go_service_broker/utils"
)

var (
	conf = config.GetConfig()
)

type Server struct {
	controller *Controller
}

func CreateServer() *Server {
	var serviceInstancesMap map[string]*module.ServiceInstance
	var keyMap map[string]*module.ServiceKey

	err := utils.ReadAndUnmarshal(&serviceInstancesMap, conf.DataPath, conf.ServiceInstancesFileName)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("WARNING: service instance data file '%s' does not exist: \n", conf.ServiceInstancesFileName)
		serviceInstancesMap = make(map[string]*module.ServiceInstance)
	}

	err = utils.ReadAndUnmarshal(&keyMap, conf.DataPath, conf.ServiceKeysFileName)
	if err != nil {
		fmt.Printf("WARNING: key map data file '%s' does not exist: \n", conf.ServiceKeysFileName)
		keyMap = make(map[string]*module.ServiceKey)
	}

	return &Server{
		controller: &Controller{
			InstanceMap: serviceInstancesMap,
			KeyMap:      keyMap,
		},
	}
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

	fmt.Println("Listening on port 8001 ..")
	http.ListenAndServe(":8001", nil)
}
