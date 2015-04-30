package model

type ServiceBinding struct {
	Id                string
	ServiceId         string
	AppId             string
	ServicePlanId     string
	PrivateKey        string
	ServiceInstanceId string
}

type CreateServiceBindingResponse struct {
	// SyslogDrainUrl string      `json:"syslog_drain_url, omitempty"`
	Credentials interface{} `json:"credentials"`
}

type Credential struct {
	UserName   string `json:"username"`
	PrivateKey string `json:"private_key"`
}
