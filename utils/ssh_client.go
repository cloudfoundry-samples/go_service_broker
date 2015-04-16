package utils

import (
	"code.google.com/p/go.crypto/ssh"
)

type SshClient interface {
	Close() error
	ExecCommand(command string) (string, error)
}

type awsSshClient struct {
	client *ssh.Client
}

func GetSshClient(username string, privateKey []byte, ip string) (*awsSshClient, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", ip+":22", config)

	return &awsSshClient{client: client}, err
}

func (sshClient *awsSshClient) Close() error {
	return sshClient.client.Close()
}

func (sshClient *awsSshClient) ExecCommand(command string) (string, error) {
	session, err := sshClient.client.NewSession()
	if err != nil {
	}
	defer session.Close()

	output, err := session.Output(command)

	return string(output), err
}
