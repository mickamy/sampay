package job

import (
	"context"

	"github.com/mickamy/slogger"

	"mickamy.com/sampay/internal/lib/aws/ses"
)

type SendEmail struct {
	sesClient ses.Client
}

func NewSendEmail(sesClient ses.Client) SendEmail {
	return SendEmail{
		sesClient: sesClient,
	}
}

func (j SendEmail) Execute(ctx context.Context, payloadStr string) error {
	// TODO: Implement the job logic
	slogger.InfoCtx(ctx, "SendEmail.Execute", "payloadStr", payloadStr)
	return nil
}
