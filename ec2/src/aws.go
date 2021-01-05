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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String("eu-west-3")},
	}))

	// Create new EC2 client
	svc := ec2.New(sess)

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
				fmt.Println(*instance.InstanceId)

				input := &ec2.TerminateInstancesInput{
					InstanceIds: []*string{
						aws.String(*instance.InstanceId),
					},
					DryRun: aws.Bool(false),
				}
				_, err := svc.TerminateInstances(input)
				awsErr, _ := err.(awserr.Error)
				if err == nil {
					fmt.Println("Success shutdown ")
				} else {
					fmt.Println("Error", awsErr)
				}
			}
		} else {
			fmt.Println("No instances running")
		}
	}
}
