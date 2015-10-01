package client

import (
	"os"
	"fmt"
	"errors"
	"strconv"

	softlayer "github.com/maximilien/softlayer-go/softlayer"
	slclient "github.com/maximilien/softlayer-go/client"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

type virtualGuestProps struct {
	hostname string
	domain string
	startCpus int
	maxMemory int
	dataCenterName string
	operatingSystemReferenceCode string
}

type SoftLayerClient struct {
	vgProps virtualGuestProps
}

func NewSoftLayerClient() *SoftLayerClient {
	fmt.Println("NewSoftLayerClient ready!")

	defatultProps := defaultVirtualGuestProperties()

	return &SoftLayerClient{
		vgProps: defatultProps,		
	}	
}

// state == pending, running, succeeded, failed
func (c *SoftLayerClient) GetInstanceState(instanceId string) (string, error) {
	vgId, err := strconv.Atoi(instanceId)
	if err != nil {
		return "failed", err
	}

	client, err := c.createSoftLayerClient()
	if err != nil {
		return "", err
	}

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
		return "failed", err
	}

	vgPowerState, err := virtualGuestService.GetPowerState(vgId)
	if err != nil {
		return "failed", err
	}

	if vgPowerState.KeyName == "RUNNING" {
		return "running", nil
	}

	return "pending", nil
}

func (c *SoftLayerClient) CreateInstance(parameters interface{}) (string, error) {
	virtualGuestTemplate := c.createVirtualGuestTemplate(parameters)

	client, err := c.createSoftLayerClient()
	if err != nil {
		return "", err
	}

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
	  return "", err
	}

	virtualGuest, err := virtualGuestService.CreateObject(virtualGuestTemplate)
	if err != nil {
	    return "", err
	}

	return strconv.Itoa(virtualGuest.Id), nil
}

func (c *SoftLayerClient) DeleteInstance(instanceId string) error {
	vgId, err := strconv.Atoi(instanceId)
	if err != nil {
		return err
	}

	client, err := c.createSoftLayerClient()
	if err != nil {
		return err
	}

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
		return err
	}

	_, err = virtualGuestService.DeleteObject(vgId)	
	if err != nil {
		return err
	}

	return nil
}

func (c *SoftLayerClient) InjectKeyPair(instanceId string) (string, string, string, error) {
	return "", "", "", nil
}

func (c *SoftLayerClient) RevokeKeyPair(instanceId string, privateKeyName string) error {
	client, err := c.createSoftLayerClient()
	if err != nil {
		return err
	}

	vgId, err := strconv.Atoi(instanceId)
	if err != nil {
		return err
	}

	sshKey, err := c.findSshKey(vgId, privateKeyName)
	if err != nil {
		return err
	}

	sshKeyService, err := client.GetSoftLayer_Security_Ssh_Key_Service()
	if err != nil {
		return err
	}

	deleted, err := sshKeyService.DeleteObject(sshKey.Id)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New(fmt.Sprintf("Could not delete ssh key with id: %d", sshKey.Id))
	}

	return nil
}

// Private methods

func (c *SoftLayerClient) createVirtualGuestTemplate(parameters interface{}) datatypes.SoftLayer_Virtual_Guest_Template {
	return datatypes.SoftLayer_Virtual_Guest_Template{
  		Hostname:  c.vgProps.hostname,
	    Domain:    c.vgProps.domain,
	    StartCpus: c.vgProps.startCpus,
	    MaxMemory: c.vgProps.maxMemory,
	    Datacenter: datatypes.Datacenter{
	        	Name: c.vgProps.dataCenterName,
	    },
	    SshKeys:                      []datatypes.SshKey{},
	    HourlyBillingFlag:            true,
	    LocalDiskFlag:                true,
	    OperatingSystemReferenceCode: c.vgProps.operatingSystemReferenceCode,
	}
}

func (c *SoftLayerClient) createSoftLayerClient() (softlayer.Client, error) {
	username := os.Getenv("SL_USERNAME")
	if username == "" {
		return nil, errors.New("You must set environment variable SL_USERNAME for SoftLayer cloud")
	}

	apiKey := os.Getenv("SL_API_KEY")
	if apiKey == "" {
		return nil, errors.New("You must set environment variable SL_API_KEY for SoftLayer cloud")
	}

	return slclient.NewSoftLayerClient(username, apiKey), nil
}

func (c *SoftLayerClient) findSshKey(vgId int, keyName string) (datatypes.SoftLayer_Security_Ssh_Key, error) {
	client, err := c.createSoftLayerClient()
	if err != nil {
		return datatypes.SoftLayer_Security_Ssh_Key{}, err
	}

	accountService, err := client.GetSoftLayer_Account_Service()
	if err != nil {
		return datatypes.SoftLayer_Security_Ssh_Key{}, err
	}

	sshKeys, err := accountService.GetSshKeys()
	if err != nil {
		return datatypes.SoftLayer_Security_Ssh_Key{}, err
	}

	for _, key := range sshKeys {
		if key.Key == keyName {
			return key, nil
		}
	}

	return datatypes.SoftLayer_Security_Ssh_Key{}, errors.New(fmt.Sprintf("Could not find ssh key with name: %s", keyName))
}

// Private functions

func defaultVirtualGuestProperties() virtualGuestProps {
	return virtualGuestProps {
		hostname: "go-service-broker",
		domain: "softlayer.com",
		startCpus: 1,
		maxMemory: 1024,
		dataCenterName: "ams01",
		operatingSystemReferenceCode: "UBUNTU_LATEST",
	}
}
