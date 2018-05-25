package sqs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/inspectr/backend/plugins"
	"github.com/spf13/viper"
)

type SQS struct {
	events chan transistor.Event
	queue  Queue
}

type SQSMessage struct {
	sqsResponse   aws_sqs.ReceiveMessageOutput
	trail         plugins.Trail
	success       bool
	statusMessage string
}

type Queue struct {
	Client sqsiface.SQSAPI
	URL    string
}

func init() {
	transistor.RegisterPlugin("sqs", func() transistor.Plugin {
		return &SQS{}
	})
}

func (x *SQS) Subscribe() []string {
	return []string{
		"trail:status",
	}
}

func (x *SQS) Process(e transistor.Event) error {
	switch e.State {
	case transistor.GetState("complete"):
		return x.deleteMessages(e)
	case transistor.GetState("failed"):
		return nil
	}
	return nil
}

func (x *SQS) Start(e chan transistor.Event) error {
	x.events = e

	region := viper.GetString("plugins.sqs.aws_region")
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					viper.GetString("plugins.sqs.aws_access_key_id"),
					viper.GetString("plugins.sqs.aws_secret_access_key"),
					"",
				),
				Region: &region,
			},
		},
	))

	queue := Queue{
		Client: aws_sqs.New(sess),
		URL:    viper.GetString("plugins.sqs.aws_sqs_url"),
	}

	x.queue = queue

	log.Info("Started SQS")

	go func(x *SQS, e chan transistor.Event) {
		for {
			msgs, err := x.queue.GetMessages(1, 20)
			if err != nil {
				log.Error(err)
			}

			for _, msg := range msgs {
				if msg.success == true {
					e <- transistor.NewEvent(transistor.EventName("trail"), transistor.GetAction("create"), msg.trail)
				} else if msg.success == false {
					x.deleteMessages(transistor.NewEvent(transistor.EventName("trail"), transistor.GetAction("delete"), msg.trail))
				}
			}

			time.Sleep(1)
		}
	}(x, e)

	return nil
}

func (x *SQS) Stop() {
	log.Info("Stopping SQS")
}

// GetMessages returns the parsed messages from SQS if any. If an error
// occurs that error will be returned.
func (q *Queue) GetMessages(numMessages int64, waitTimeout int64) ([]SQSMessage, error) {
	params := aws_sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.URL),
		MaxNumberOfMessages: aws.Int64(numMessages),
	}
	if waitTimeout > 0 {
		params.WaitTimeSeconds = aws.Int64(waitTimeout)
	}
	resp, err := q.Client.ReceiveMessage(&params)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages, %v", err)
	}

	log.Info("Succesfully got SQS Messages")

	msgs := make([]SQSMessage, len(resp.Messages))
	for i, msg := range resp.Messages {

		sqsMessage := SQSMessage{}
		parsedTrail := plugins.Trail{}
		parsedTrail.MessageID = *msg.ReceiptHandle
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedTrail); err != nil {
			sqsMessage.success = false
			sqsMessage.statusMessage = fmt.Sprintf("failed to unmarshal message, %v", err)
		}
		sqsMessage.success = true

		sqsMessage.trail = parsedTrail

		msgs[i] = sqsMessage
	}

	return msgs, nil
}

func (x *SQS) deleteMessages(e transistor.Event) error {
	if e.PayloadModel == "plugins.Trail" {

		msg := e.Payload.(plugins.Trail)

		deleteParams := aws_sqs.DeleteMessageInput{
			QueueUrl:      aws.String(x.queue.URL),
			ReceiptHandle: &msg.MessageID,
		}

		_, err := x.queue.Client.DeleteMessage(&deleteParams)
		if err != nil {
			log.Info("Failed to delete SQS Message")
			return err
		} else {
			log.Info("Deleted SQS Message")
			return nil
		}

	}
	return nil
}
