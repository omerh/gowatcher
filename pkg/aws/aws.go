package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"log"
	"os"
)


var mySession *session.Session
var sessionAutoScaling *autoscaling.AutoScaling

func initSession(instanceRegion string) (*session.Session,error){
	sessionConf := aws.NewConfig().WithRegion(instanceRegion)
	thisSession, _ := session.NewSession(sessionConf)
	mySession = thisSession
	return mySession, nil
}

func initAutoScalingAwsSession(awsSession *session.Session) (*autoscaling.AutoScaling, error) {
	sessionAutoScaling = autoscaling.New(awsSession)
	return sessionAutoScaling, nil
}

///TerminateInstance function start termination on a ec2 instance
func TerminateInstance() {
	log.Println("Starting to run terminate instance process")
	// Get instance id and region from metadata
	instanceId, instanceRegion := getInstanceID()
	log.Printf("Working on %v in %v region", instanceId, instanceRegion)

	// Init aws session
	awsSession,_ := initSession(instanceRegion)
	log.Println("Initialized aws session")

	// Init Aws auto scaling session
	initAutoScalingAwsSession(awsSession)
	log.Println("Initialized auto scaling session")

	// Get auto scaling group name
	instanceAutoScaleGroupName := getAutoScalingName(instanceId)
	log.Printf("Instance %v auto scaling group name is: %v", instanceId, instanceAutoScaleGroupName)

	// Set instance scale in policy to false
	success := setScaleInProtectionToInstance(instanceAutoScaleGroupName, instanceId)

	// Terminate ec2 instance after setting scale in policy to false
	if success{
		terminateInstance(instanceId)
	}
}

func terminateInstance(instanceId string) bool {
	_, ok := os.LookupEnv("DEBUG")
	if !ok {
		// input filter with instanceId of its auto scaling group
		input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
			InstanceId:                     aws.String(instanceId),
			ShouldDecrementDesiredCapacity: aws.Bool(true),
		}

		// Terminate instance from auto scaling group
		result, err := sessionAutoScaling.TerminateInstanceInAutoScalingGroup(input)

		if err != nil {
			log.Println(err)
			return false
		} else {
			log.Println(result.String())
			return true
		}
	} else {
		log.Printf("In Debug mode, terminate to instance %v should be terminated", instanceId)
		return true
	}
}

func setScaleInProtectionToInstance(instanceAutoScaleGroupName string, instanceId string) bool {
	// Set input filter of instance id and auto scale group for setting protect false
	input := &autoscaling.SetInstanceProtectionInput{
		AutoScalingGroupName: aws.String(instanceAutoScaleGroupName),
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		ProtectedFromScaleIn: aws.Bool(false),
	}

	// Set Instance protection on auto scale group
	result, err := sessionAutoScaling.SetInstanceProtection(input)

	if err !=  nil {
		log.Panic(err)
		return false
	} else {
		log.Printf("Instance %v autoscale protect from scaleIn was set for group %v", instanceId, instanceAutoScaleGroupName)
		log.Println(result.String())
		return true
	}
}

func getAutoScalingName(instanceId string) string {
	input := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	// Locate Auto scale group for instance
	ec2AutoScaleInstanceDetails, err := sessionAutoScaling.DescribeAutoScalingInstances(input)

	if err != nil {
		log.Println("Problem running DescribeAutoScalingInstances")
		panic(err)
	}

	// Extract AutoScale group name
	autoScaleGroupName := *(ec2AutoScaleInstanceDetails.AutoScalingInstances[0].AutoScalingGroupName)
	log.Printf("AutoScaling group name for instance %v is %v", instanceId, autoScaleGroupName)

	//return autoScaleSlice
	return autoScaleGroupName
}

func getInstanceID() (instanceId string, instanceRegion string) {
	// Load session details from instance metadata
	svc := ec2metadata.New(session.New())

	// Get from metadata instance Id and region
	ec2Instance, err := svc.GetInstanceIdentityDocument()

	if err != nil {
		panic(err)
	}

	log.Printf("Instance id is: %v, Instance Regions is %v", ec2Instance.InstanceID, ec2Instance.Region)
	return ec2Instance.InstanceID, ec2Instance.Region
}