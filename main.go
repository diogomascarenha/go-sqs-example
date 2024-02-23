package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gookit/goutil/dump"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	awsSQSAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSQSSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSQSSessionToken := os.Getenv("AWS_SESSION_TOKEN")

	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			awsSQSAccessKeyId,
			awsSQSSecretKey,
			awsSQSSessionToken,
		),
	})

	svc := sqs.New(sess)

	listSQSQueuesExample(svc)
	sendSQSMessageExample(svc)
	retrieveSQSMessageExample(svc)

}

func listSQSQueuesExample(svc *sqs.SQS) {
	result, err := svc.ListQueues(nil)
	if err != nil {
		panic(err)
	}
	for i, url := range result.QueueUrls {
		fmt.Printf("%d: %s\n", i, *url)
	}
}

func sendSQSMessageExample(svc *sqs.SQS) {
	messageGroupId := aws.String("group-id=grupo-diogo-teste-01")
	messageDeduplicationId := aws.String("258e7707-18fb-4006-8f3a-9ff0fa0d0780")
	queueUrl := aws.String(os.Getenv("AWS_SQS_QUEUE_URL"))
	messageBody := aws.String("Message Body Example")
	messageAttributes := map[string]*sqs.MessageAttributeValue{
		"Title": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("Title Example"),
		},
		"Author": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("Author Example"),
		},
	}

	sendMessageOutput, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         messageGroupId,
		MessageDeduplicationId: messageDeduplicationId,
		QueueUrl:               queueUrl,
		MessageBody:            messageBody,
		//DelaySeconds: aws.Int64(10),
		MessageAttributes: messageAttributes,
	})

	if err != nil {
		panic(err)
	}

	dump.Println(sendMessageOutput)
}

func retrieveSQSMessageExample(svc *sqs.SQS) {
	timeout := flag.Int64("t", 5, "How long, in seconds, that the message is hidden from others")
	flag.Parse()

	if *timeout < 0 {
		*timeout = 0
	}

	if *timeout > 12*60*60 {
		*timeout = 12 * 60 * 60
	}

	queueUrl := aws.String(os.Getenv("AWS_SQS_QUEUE_URL"))

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   timeout,
	})

	if err != nil {
		panic(err)
	}

	dump.Println(*msgResult)
}
