package scenario_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
)

func TestEmailVerification(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	infras := di.NewInfras(newReadWriter(t), newKVS(t))
	server := initServer(t, infras)

	token := emailVerification(t, ctx, infras, server)
	assert.NotEmpty(t, token)
}

func requestEmailVerification(t *testing.T, s *httptest.Server, intentType authModel.EmailVerificationIntentType, email string, f func(res *connect.Response[authv1.RequestVerificationResponse], err error)) {
	t.Helper()

	var intent authv1.RequestVerificationRequest_IntentType
	switch intentType {
	case authModel.EmailVerificationIntentTypeSignUp:
		intent = authv1.RequestVerificationRequest_INTENT_TYPE_SIGN_UP
	case authModel.EmailVerificationIntentTypeResetPassword:
		intent = authv1.RequestVerificationRequest_INTENT_TYPE_RESET_PASSWORD
	default:
		t.Fatalf("unexpected intent type: %v", intentType)
	}
	client := authv1connect.NewEmailVerificationServiceClient(http.DefaultClient, s.URL+"/api")
	req := connect.NewRequest(&authv1.RequestVerificationRequest{
		IntentType: intent,
		Email:      email,
	})
	res, err := client.RequestVerification(context.Background(), req)
	f(res, err)
}

func verifyEmail(t *testing.T, s *httptest.Server, requestToken string, pinCode string, f func(res *connect.Response[authv1.VerifyEmailResponse], err error)) {
	t.Helper()

	client := authv1connect.NewEmailVerificationServiceClient(http.DefaultClient, s.URL+"/api")
	req := connect.NewRequest(&authv1.VerifyEmailRequest{
		Token:   requestToken,
		PinCode: pinCode,
	})
	res, err := client.VerifyEmail(context.Background(), req)
	f(res, err)
}

func emailVerification(t *testing.T, ctx context.Context, infras di.Infras, s *httptest.Server) string {
	t.Helper()

	var requestToken string
	{
		requestEmailVerification(t, s, authModel.EmailVerificationIntentTypeSignUp, email, func(res *connect.Response[authv1.RequestVerificationResponse], err error) {
			require.NoError(t, err)
			require.NotEmpty(t, res.Msg.Token)
			requestToken = res.Msg.Token
		})
	}

	var verifyToken string
	{
		var pinCode string
		require.NoError(t, infras.WithContext(ctx).Model(&authModel.RequestedEmailVerification{}).Where("token = ?", requestToken).Pluck("pin_code", &pinCode).Error)
		verifyEmail(t, s, requestToken, pinCode, func(res *connect.Response[authv1.VerifyEmailResponse], err error) {
			require.NoError(t, err)
			require.NotEmpty(t, res.Msg.Token)
			verifyToken = res.Msg.Token
		})
	}

	t.Logf("verifyToken: %s", verifyToken)
	return verifyToken
}
