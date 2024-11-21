package main

import (
	"context"
	"encoding/json"
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

		dataEvent := map[string]interface{}{
			"data": map[string]string{"message": messageBody},
		}

		// Converta o evento em um formato JSON válido
		dataEventBytes, err := json.Marshal(dataEvent)
		if err != nil {
			log.Fatal("Erro ao criar evento:", err)
		}

		event := &eventbridge.PutEventsRequestEntry{
			DetailType:   aws.String("SQSMessage"),
			Source:       aws.String("my.lambda"),
			Detail:       aws.String(string(dataEventBytes)),
			EventBusName: aws.String("LambdaHelloBus"),
		}

		log.Printf("Enviando evento para EventBridge: %v", event)

		output, err := svc.PutEvents(&eventbridge.PutEventsInput{
			Entries: []*eventbridge.PutEventsRequestEntry{event},
		})
		if err != nil {
			log.Printf("Erro ao enviar evento para EventBridge: %v", err)
			return err
		}

		log.Printf("Resultado do PutEvents: %+v", output)

		for _, entry := range output.Entries {
			if entry.ErrorCode != nil {
				log.Printf("Erro ao enviar evento: Code: %s, Message: %s",
					aws.StringValue(entry.ErrorCode), aws.StringValue(entry.ErrorMessage))
			} else {
				fmt.Printf("Evento enviado com sucesso! EventId: %s\n", aws.StringValue(entry.EventId))
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
