package web_server

import (
	"fmt"
	"net/http"

	ac "github.com/xingzhou/go_service_broker/aws_client"

	"github.com/xingzhou/go_service_broker/model"
	"github.com/xingzhou/go_service_broker/utils"
)

const (
	DEFAULT_POLLING_INTERVAL_SECONDS = 10
)

type Controller struct {
	serviceClient ac.Client

	instanceMap map[string]*model.ServiceInstance
	bindingMap  map[string]*model.ServiceBinding
}

func CreateController(instanceMap map[string]*model.ServiceInstance, bindingMap map[string]*model.ServiceBinding) *Controller {
	return &Controller{
		instanceMap:   instanceMap,
		bindingMap:    bindingMap,
		serviceClient: ac.NewClient("us-east-1"),
	}
}

func (c *Controller) Catalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Service Broker Catalog...")

	var catalog model.Catalog
	err := utils.ReadAndUnmarshal(&catalog, conf.CatalogPath, "catalog.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, http.StatusOK, catalog)
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Service Instance...")

	var instance model.ServiceInstance

	err := utils.ProvisionDataFromRequest(r, &instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instanceId, err := c.serviceClient.CreateInstance(instance.Parameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instance.InternalId = instanceId
	instance.DashboardUrl = "http://dashbaord_url"
	instance.Id = utils.ExtractVarsFromRequest(r, "service_instance_guid")
	instance.LastOperation = &model.LastOperation{
		State:                    "in progress",
		Description:              "creating service instance...",
		AsyncPollIntervalSeconds: DEFAULT_POLLING_INTERVAL_SECONDS,
	}

	c.instanceMap[instance.Id] = &instance
	err = utils.MarshalAndRecord(c.instanceMap, conf.DataPath, conf.ServiceInstancesFileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := model.CreateServiceInstanceResponse{
		DashboardUrl:  instance.DashboardUrl,
		LastOperation: instance.LastOperation,
	}
	utils.WriteResponse(w, http.StatusAccepted, response)
}

func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Service Instance State....")

	instanceId := utils.ExtractVarsFromRequest(r, "service_instance_guid")
	instance := c.instanceMap[instanceId]
	if instance == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	state, err := c.serviceClient.GetInstanceState(instance.InternalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if state == "pending" {
		instance.LastOperation.State = "in progress"
		instance.LastOperation.Description = "creating service instance..."
	} else if state == "running" {
		instance.LastOperation.State = "succeeded"
		instance.LastOperation.Description = "successfully created service instance"
	} else {
		instance.LastOperation.State = "failed"
		instance.LastOperation.Description = "failed to create service instance"
	}

	response := model.CreateServiceInstanceResponse{
		DashboardUrl:  instance.DashboardUrl,
		LastOperation: instance.LastOperation,
	}
	utils.WriteResponse(w, http.StatusOK, response)
}

func (c *Controller) RemoveServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Remove Service Instance...")

	instanceId := utils.ExtractVarsFromRequest(r, "service_instance_guid")
	instance := c.instanceMap[instanceId]
	if instance == nil {
		w.WriteHeader(http.StatusGone)
		return
	}

	err := c.serviceClient.DeleteInstance(instance.InternalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	delete(c.instanceMap, instanceId)
	utils.MarshalAndRecord(c.instanceMap, conf.DataPath, conf.ServiceInstancesFileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.deleteAssociatedBindings(instanceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, http.StatusOK, "{}")
}

func (c *Controller) Bind(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Bind Service Instance...")

	bindingId := utils.ExtractVarsFromRequest(r, "service_binding_guid")
	instanceId := utils.ExtractVarsFromRequest(r, "service_instance_guid")

	instance := c.instanceMap[instanceId]
	if instance == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ip, userName, privateKey, err := c.serviceClient.InjectKeyPair(instance.InternalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	credential := model.Credential{
		PublicIP:   ip,
		UserName:   userName,
		PrivateKey: privateKey,
	}

	response := model.CreateServiceBindingResponse{
		Credentials: credential,
	}

	c.bindingMap[bindingId] = &model.ServiceBinding{
		Id:                bindingId,
		ServiceId:         instance.ServiceId,
		ServicePlanId:     instance.PlanId,
		PrivateKey:        privateKey,
		ServiceInstanceId: instance.Id,
	}

	err = utils.MarshalAndRecord(c.bindingMap, conf.DataPath, conf.ServiceBindingsFileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, http.StatusCreated, response)
}

func (c *Controller) UnBind(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Unbind Service Instance...")

	bindingId := utils.ExtractVarsFromRequest(r, "service_binding_guid")
	instanceId := utils.ExtractVarsFromRequest(r, "service_instance_guid")
	instance := c.instanceMap[instanceId]
	if instance == nil {
		w.WriteHeader(http.StatusGone)
		return
	}

	err := c.serviceClient.RevokeKeyPair(instance.InternalId, c.bindingMap[bindingId].PrivateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	delete(c.bindingMap, bindingId)
	err = utils.MarshalAndRecord(c.bindingMap, conf.DataPath, conf.ServiceBindingsFileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, http.StatusOK, "{}")
}

func (c *Controller) deleteAssociatedBindings(instanceId string) error {
	for id, binding := range c.bindingMap {
		if binding.ServiceInstanceId == instanceId {
			delete(c.bindingMap, id)
		}
	}

	return utils.MarshalAndRecord(c.bindingMap, conf.DataPath, conf.ServiceBindingsFileName)
}
