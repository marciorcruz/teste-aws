package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type EventBridgeEvent struct {
	DetailType string `json:"detail-type"`
	Source     string `json:"source"`
	Detail     string `json:"detail"`
}

var svc *eventbridge.EventBridge

func init() {
	sess := session.Must(session.NewSession())
	svc = eventbridge.New(sess)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		messageBody := record.Body

		event := &eventbridge.PutEventsRequestEntry{
			DetailType:   aws.String("SQSMessage"),
			Source:       aws.String("my.sqs.lambda"),
			Detail:       aws.String(messageBody),
			EventBusName: aws.String("LambdaHelloBus"),
		}

		output, err := svc.PutEvents(&eventbridge.PutEventsInput{
			Entries: []*eventbridge.PutEventsRequestEntry{event},
		})
		if err != nil {
			log.Printf("Erro ao enviar evento para EventBridge: %v", err)
			return err
		}

		for _, success := range output.Entries {
			fmt.Printf("Evento enviado com sucesso! EventId: %s\n", aws.StringValue(success.EventId))
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
