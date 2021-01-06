package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	// Load session from shared config
	svc := createClient()

	// Specify the details of the instance that you want to create.
	createInstance(svc)

	reservations := listInstances(svc)

	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			terminateInstance(instance, svc)
		}
	}
}

func createInstance(svc *ec2.EC2) {
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the eu-west-3 region
		ImageId:      aws.String("ami-00798d7180f25aac2"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}

	_, error := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{{
			Key:   aws.String("env"),
			Value: aws.String("test"),
		}}})

	if error != nil {
		fmt.Println("Could not create tags", err)
		return
	}
}

func listInstances(svc *ec2.EC2) []*ec2.Reservation {

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending"), aws.String("stopped")},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	for _, reservation := range resp.Reservations {
		if len(reservation.Instances) != 0 {
			for _, instance := range reservation.Instances {
				fmt.Println(instance)
			}
		} else {
			fmt.Println("No instances running")
		}
	}

	return resp.Reservations
}

func terminateInstance(instance *ec2.Instance, svc *ec2.EC2) {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(*instance.InstanceId),
		},
		DryRun: aws.Bool(false),
	}
	_, err := svc.TerminateInstances(input)
	awsErr, _ := err.(awserr.Error)
	if err == nil {
		fmt.Println("Success shutdown")
	} else {
		fmt.Println("Error", awsErr)
	}
}

func createClient() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String("eu-west-3")},
	}))

	// Create new EC2 client
	svc := ec2.New(sess)
	return svc
}
