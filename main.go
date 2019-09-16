package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/eiladin/policy-publisher/config"
)

func main() {
	config := config.InitConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Sqs.Region),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	svc := sqs.New(sess)

	queue, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("policy-collector-queue"),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Publishing to: ", *queue.QueueUrl)

	var i rune
	for {
		fmt.Println("p to publish, q to quit")
		fmt.Scanf("%c\n", &i)
		if i == 'q' {
			break
		} else if i == 'p' {
			publishMessage(queue.QueueUrl, svc)
		}
	}
}

func publishMessage(queueUrl *string, svc *sqs.SQS) {
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(0),
		MessageAttributes: buildMessage("portal", "licensing", "v1"),
		MessageBody:       aws.String("http://localhost:5001/my/fake/service"),
		QueueUrl:          queueUrl,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Message Sent", *result.MessageId)
}

func buildMessage(product, app, version string) map[string]*sqs.MessageAttributeValue {
	result := map[string]*sqs.MessageAttributeValue{
		"Product": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(product),
		},
		"App": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(app),
		},
		"Version": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(version),
		},
	}

	return result
}
