package web_server

import (
	"os"
	"fmt"
	"net/http"
	"errors"

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

func CreateServer(cloudName string) (*Server, error) {
	serviceInstances, err := loadServiceInstances()
	if err != nil {
		return nil, err
	}

	serviceBindings, err := loadServiceBindings()
	if err != nil {
		return nil, err
	}

	controller, err := CreateController(cloudName, serviceInstances, serviceBindings)
	if err != nil {
		return nil, err
	}

	return &Server{
		controller: controller,
	}, nil
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

	cfPort := os.Getenv("PORT")
	if cfPort != "" {
		conf.Port = cfPort
	}

	fmt.Println("Server started, listening on port " + conf.Port + "...")
	fmt.Println("CTL-C to break out of broker")
	http.ListenAndServe(":"+conf.Port, nil)
}

// private methods
func loadServiceInstances() (map[string]*model.ServiceInstance, error) {
	var serviceInstancesMap map[string]*model.ServiceInstance

	err := utils.ReadAndUnmarshal(&serviceInstancesMap, conf.DataPath, conf.ServiceInstancesFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("WARNING: service instance data file '%s' does not exist: \n", conf.ServiceInstancesFileName)
			serviceInstancesMap = make(map[string]*model.ServiceInstance)
		} else {
			return nil, errors.New(fmt.Sprintf("Could not load the service instances, message: %s", err.Error()))
		}
	}

	return serviceInstancesMap, nil
}

func loadServiceBindings() (map[string]*model.ServiceBinding, error) {
	var bindingMap map[string]*model.ServiceBinding

	err := utils.ReadAndUnmarshal(&bindingMap, conf.DataPath, conf.ServiceBindingsFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("WARNING: key map data file '%s' does not exist: \n", conf.ServiceBindingsFileName)
			bindingMap = make(map[string]*model.ServiceBinding)
		} else {
			return nil, errors.New(fmt.Sprintf("Could not load the service instances, message: %s", err.Error()))
		}
	}

	return bindingMap, nil
}
