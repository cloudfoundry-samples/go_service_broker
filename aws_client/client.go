package aws_client

import (
	"fmt"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

type Client interface {
	CreateInstance(imageId string) error
}

type AWSClient struct {
	Region string
}

func NewClient(region string) *AWSClient {
	return &AWSClient{
		Region: region,
	}
}

func (c *AWSClient) CreateInstance() error {
	fmt.Println("creating instance...")
	return c.createInstance("ami-80280ee8")
}

func (c *AWSClient) createInstance(imageId string) error {
	// Create an EC2 service object in the "us-west-2" region
	// Note that you can also configure your region globally by
	// exporting the AWS_REGION environment variable
	svc := ec2.New(&aws.Config{Region: c.Region})
	fmt.Println("EC2 service initialized!")

	// Call the DescribeInstances Operation
	// resp, err := svc.DescribeInstances(nil)
	// fmt.Println("resp:", resp)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	panic(err)
	// }

	// // resp has all of the response data, pull out instance IDs:
	// fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
	// for idx, res := range resp.Reservations {
	// 	fmt.Println("  > Number of instances: ", len(res.Instances))
	// 	for _, inst := range resp.Reservations[idx].Instances {
	// 		fmt.Println("    - Instance ID: ", *inst.InstanceID)
	// 	}
	// }

	// input := &ec2.DescribeImagesInput{
	// 	// DryRun: aws.Boolean(true),
	// 	// ImageIDs: []*string{
	// 	// 	aws.String("ami-01da0968"),
	// 	// },
	// 	Filters: []*ec2.Filter{
	// 		&ec2.Filter{
	// 			Name: aws.String("image-id"),
	// 			Values: []*string{
	// 				aws.String(imageId),
	// 			},
	// 		},
	// 		// &ec2.Filter{
	// 		// 	Name: aws.String("is-public"),
	// 		// 	Values: []*string{
	// 		// 		aws.String("true"),
	// 		// 	},
	// 		// },
	// 	},
	// }

	// // output := &ec2.DescribeImagesOutput{}
	// fmt.Println("Describing images...")
	// output, err := svc.DescribeImages(input)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	panic(err)
	// }

	// fmt.Println(awsutil.StringValue(output))

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
		InstanceType: aws.String("t2.micro"),
		// KernelID:                          aws.String("String"),
		// KeyName:                           aws.String("String"),
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
			aws.String("sg-b23aead6"), // Required
			// More values...
		},
		// SecurityGroups: []*string{
		// 	aws.String("default"), // Required
		// 	// More values...
		// },
		SubnetID: aws.String("subnet-0c75a427"),
		UserData: aws.String("service async key demo"),
	}

	instancOutput, err := svc.RunInstances(instanceInput)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(instancOutput))
	return nil
}
