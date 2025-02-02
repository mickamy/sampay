package job

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"mickamy.com/sampay/internal/lib/aws/ses"
)

type SendEmailPayload struct {
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type SendEmail struct {
	sesClient ses.Client
}

func NewSendEmail(sesClient ses.Client) SendEmail {
	return SendEmail{
		sesClient: sesClient,
	}
}

func (j SendEmail) Execute(ctx context.Context, payloadStr string) error {
	var payload SendEmailPayload
	if err := parsePayload(ctx, payloadStr, &payload); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("\"Sampay\" <" + payload.From + ">"),
		Destination: &types.Destination{
			ToAddresses: []string{payload.To},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(payload.Subject),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(payload.Body),
					},
				},
			},
		},
	}

	_, err := j.sesClient.Send(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
