package aws_client

import (
	"fmt"
	"strconv"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/ec2"

	"github.com/xingzhou/go_service_broker/utils"
)

const (
	AMI_ID            = "ami-dc5e75b4" //"ami-ecb68a84"
	SECURITY_GROUP_ID = "sg-b23aead6"
	SUBNET_ID         = "subnet-0c75a427"
	KEY_NAME          = "mykey1"
	INSTANCE_TYPE     = "t2.micro"
	LINUX_USER        = "ubuntu"
	PRIVATE_KEY_PATH  = "/Users/dongdong/.ssh/id_rsa"
)

type Client interface {
	CreateInstance() (string, error)
	GetInstanceState(instanceId string) (string, error)
	CreateKeyPair(keyName string) (string, error)
	InjectKeyPair(instanceId string) (string, error)
}

type AWSClient struct {
	EC2Client *ec2.EC2
}

func NewClient(region string) *AWSClient {
	return &AWSClient{
		EC2Client: ec2.New(&aws.Config{Region: region}),
	}
}

func (c *AWSClient) CreateInstance() (string, error) {
	return c.createInstance(AMI_ID)
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

func (c *AWSClient) CreateKeyPair(keyName string) (string, error) {
	keypairInput := &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	}

	keypairOutput, err := c.EC2Client.CreateKeyPair(keypairInput)
	if err != nil {
		return "", err
	}
	return awsutil.StringValue(keypairOutput), nil
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
	pemBytes, err := utils.ReadFile(PRIVATE_KEY_PATH)
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
	cat ./broker_id_rsa
	`
	privateKey, err := awsSShClient.ExecCommand(command)
	if err != nil {
		return "", err
	}

	return privateKey, nil
}

func (c *AWSClient) createInstance(imageId string) (string, error) {
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
		KeyName: aws.String(KEY_NAME),
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
