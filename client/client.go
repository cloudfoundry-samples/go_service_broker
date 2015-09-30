package client

type Client interface {
	CreateInstance(parameters interface{}) (string, error)
	GetInstanceState(instanceId string) (string, error)
	InjectKeyPair(instanceId string) (string, string, string, error)
	DeleteInstance(instanceId string) error
	RevokeKeyPair(instanceId string, privateKey string) error
}