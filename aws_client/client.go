package aws_client

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/ec2"

	"github.com/xingzhou/go_service_broker/utils"
)

const (
	AMI_ID                = "ami-dc5e75b4" //"ami-ecb68a84"
	SECURITY_GROUP_ID     = "sg-b23aead6"
	SUBNET_ID             = "subnet-0c75a427"
	KEYPAIR_NAME          = "broker_keypair"
	INSTANCE_TYPE         = "t2.micro"
	LINUX_USER            = "ubuntu"
	KEYPAIR_DIR_NAME      = ".gsb"
	PIRVATE_KEY_FILE_NAME = "broker_id_rsa"
	PUBLIC_KEY_FILE_NAME  = "broker_id_rsa.pub"
)

type Client interface {
	CreateInstance(parameters interface{}) (string, error)
	GetInstanceState(instanceId string) (string, error)
	InjectKeyPair(instanceId string) (string, error)
	DeleteInstance(instanceId string) error
	RevokeKeyPair(instanceId string, privateKey string) error
}

type AWSClient struct {
	EC2Client *ec2.EC2
}

func NewClient(region string) *AWSClient {
	return &AWSClient{
		EC2Client: ec2.New(&aws.Config{Region: region}),
	}
}

func (c *AWSClient) CreateInstance(parameters interface{}) (string, error) {
	var amiId string

	switch parameters.(type) {
	case map[string]interface{}:
		param := parameters.(map[string]interface{})
		amiId = param["ami_id"].(string)
	default:
		amiId = AMI_ID
	}

	return c.createInstance(amiId)
}

func (c *AWSClient) GetInstanceState(instanceId string) (string, error) {
	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIDs: []*string{
			aws.String(instanceId), // Required
		},
	}

	instanceOutput, err := c.EC2Client.DescribeInstances(instanceInput)
	if err != nil {
		return "", err
	}

	state, _ := strconv.Unquote(awsutil.StringValue(instanceOutput.Reservations[0].Instances[0].State.Name))
	return state, nil
}

func (c *AWSClient) setupKeyPair() error {
	private_key_file := path.Join(os.Getenv("HOME"), KEYPAIR_DIR_NAME, PIRVATE_KEY_FILE_NAME)

	if !utils.Exists(private_key_file) {
		keypairInput := &ec2.CreateKeyPairInput{
			KeyName: aws.String(KEYPAIR_NAME),
		}

		keypairOutput, err := c.EC2Client.CreateKeyPair(keypairInput)
		if err != nil {
			return err
		}

		key_dir := path.Join(os.Getenv("HOME"), KEYPAIR_DIR_NAME)
		if !utils.MkDir(key_dir) {
			return errors.New("failed to create local keypair directory")
		}

		key_data, _ := strconv.Unquote(awsutil.StringValue(keypairOutput.KeyMaterial))
		err = utils.WriteFile(private_key_file, []byte(key_data))
		if err != nil {
			return errors.New("failed to save private key file")
		}
	}

	return nil
}

func (c *AWSClient) InjectKeyPair(instanceId string) (string, error) {
	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIDs: []*string{
			aws.String(instanceId), // Required
		},
	}

	instanceOutput, err := c.EC2Client.DescribeInstances(instanceInput)
	if err != nil {
		return "", err
	}

	ip, _ := strconv.Unquote(awsutil.StringValue(instanceOutput.Reservations[0].Instances[0].PublicIPAddress))
	pemBytes, err := utils.ReadFile(path.Join(os.Getenv("HOME"), KEYPAIR_DIR_NAME, PIRVATE_KEY_FILE_NAME))
	if err != nil {
		return "", err
	}

	awsSShClient, err := utils.GetSshClient(LINUX_USER, pemBytes, ip)
	if err != nil {
		return "", err
	}

	command := `rm -f ./broker_id_rsa ./broker_id_rsa.pub
		ssh-keygen -q -t rsa -N ""  -f ./broker_id_rsa
		cat ./broker_id_rsa.pub >> .ssh/authorized_keys
		cat ./broker_id_rsa`

	privateKey, err := awsSShClient.ExecCommand(command)
	if err != nil {
		return "", err
	}

	return privateKey, nil
}

func (c *AWSClient) createInstance(imageId string) (string, error) {
	err := c.setupKeyPair()
	if err != nil {
		return "", err
	}

	instanceInput := &ec2.RunInstancesInput{
		ImageID:  aws.String(imageId), // Required
		MaxCount: aws.Long(1),         // Required
		MinCount: aws.Long(1),         // Required
		// AdditionalInfo: aws.String("String"),
		// BlockDeviceMappings: []*ec2.BlockDeviceMapping{
		// 	&ec2.BlockDeviceMapping{ // Required
		// 		DeviceName: aws.String("String"),
		// 		EBS: &ec2.EBSBlockDevice{
		// 			DeleteOnTermination: aws.Boolean(true),
		// 			Encrypted:           aws.Boolean(true),
		// 			IOPS:                aws.Long(1),
		// 			SnapshotID:          aws.String("String"),
		// 			VolumeSize:          aws.Long(1),
		// 			VolumeType:          aws.String("VolumeType"),
		// 		},
		// 		NoDevice:    aws.String("String"),
		// 		VirtualName: aws.String("String"),
		// 	},
		// 	// More values...
		// },
		// ClientToken: aws.String("String"),
		// DisableAPITermination: aws.Boolean(true),
		// DryRun:                aws.Boolean(true),
		// EBSOptimized:          aws.Boolean(true),
		// IAMInstanceProfile: &ec2.IAMInstanceProfileSpecification{
		// 	ARN:  aws.String("String"),
		// 	Name: aws.String("String"),
		// },
		// InstanceInitiatedShutdownBehavior: aws.String("ShutdownBehavior"),
		InstanceType: aws.String(INSTANCE_TYPE),
		// KernelID:                          aws.String("String"),
		KeyName: aws.String(KEYPAIR_NAME),
		// Monitoring: &ec2.RunInstancesMonitoringEnabled{
		// 	Enabled: aws.Boolean(true), // Required
		// },
		// NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
		// 	&ec2.InstanceNetworkInterfaceSpecification{ // Required
		// 		AssociatePublicIPAddress: aws.Boolean(true),
		// 		DeleteOnTermination:      aws.Boolean(true),
		// 		Description:              aws.String("String"),
		// 		DeviceIndex:              aws.Long(1),
		// 		Groups: []*string{
		// 			aws.String("String"), // Required
		// 			// More values...
		// 		},
		// 		NetworkInterfaceID: aws.String("String"),
		// 		PrivateIPAddress:   aws.String("String"),
		// 		PrivateIPAddresses: []*ec2.PrivateIPAddressSpecification{
		// 			&ec2.PrivateIPAddressSpecification{ // Required
		// 				PrivateIPAddress: aws.String("String"), // Required
		// 				Primary:          aws.Boolean(true),
		// 			},
		// 			// More values...
		// 		},
		// 		SecondaryPrivateIPAddressCount: aws.Long(1),
		// 		SubnetID:                       aws.String("String"),
		// 	},
		// 	// More values...
		// },
		// Placement: &ec2.Placement{
		// 	AvailabilityZone: aws.String("String"),
		// 	GroupName:        aws.String("String"),
		// 	Tenancy:          aws.String("Tenancy"),
		// },
		// PrivateIPAddress: aws.String("String"),
		// RAMDiskID:        aws.String("String"),
		SecurityGroupIDs: []*string{
			aws.String(SECURITY_GROUP_ID), // Required
			// More values...
		},
		SubnetID: aws.String(SUBNET_ID),
	}

	instanceOutput, err := c.EC2Client.RunInstances(instanceInput)
	if err != nil {
		return "", err
	}

	fmt.Println(awsutil.StringValue(instanceOutput))
	instanceId, _ := strconv.Unquote(awsutil.StringValue(instanceOutput.Instances[0].InstanceID))

	return instanceId, nil
}

func (c *AWSClient) DeleteInstance(instanceId string) error {
	terminateInstanceInput := &ec2.TerminateInstancesInput{
		// One or more instance IDs.
		InstanceIDs: []*string{
			aws.String(instanceId), // Required
		},
	}

	_, err := c.EC2Client.TerminateInstances(terminateInstanceInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *AWSClient) RevokeKeyPair(instanceId string, privateKey string) error {
	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIDs: []*string{
			aws.String(instanceId),
		},
	}

	instanceOutput, err := c.EC2Client.DescribeInstances(instanceInput)
	if err != nil {
		return err
	}

	ip, _ := strconv.Unquote(awsutil.StringValue(instanceOutput.Reservations[0].Instances[0].PublicIPAddress))
	pemBytes, err := utils.ReadFile(path.Join(os.Getenv("HOME"), KEYPAIR_DIR_NAME, PIRVATE_KEY_FILE_NAME))
	if err != nil {
		return err
	}

	awsSShClient, err := utils.GetSshClient(LINUX_USER, pemBytes, ip)
	if err != nil {
		return err
	}

	publicKey, err := utils.GeneratePublicKey([]byte(privateKey))
	if err != nil {
		return err
	}

	escapedPublicKey := strings.Replace(publicKey, "/", "\\/", -1)
	command := fmt.Sprintf("sed '/%s/d' -i ~/.ssh/authorized_keys && echo 'revoked the public key: %s'", escapedPublicKey, publicKey)

	result, err := awsSShClient.ExecCommand(command)
	if err != nil {
		return err
	}
	fmt.Println(result)

	return nil
}
