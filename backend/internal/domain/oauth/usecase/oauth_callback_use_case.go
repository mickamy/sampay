package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	oauthModel "mickamy.com/sampay/internal/domain/oauth/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/aws/s3"
	"mickamy.com/sampay/internal/lib/oauth"
)

type OAuthCallbackInput struct {
	Provider oauthModel.OAuthProvider
	Code     string
}

type OAuthCallbackOutput struct {
	Session authModel.Session
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type OAuthCallback interface {
	Do(ctx context.Context, input OAuthCallbackInput) (OAuthCallbackOutput, error)
}

type oauthCallback struct {
	google                oauth.Google
	s3                    s3.Client
	writer                *database.Writer
	authRepo              authRepository.Authentication
	emailVerificationRepo authRepository.EmailVerification
	sessionRepo           authRepository.Session
	userRepo              userRepository.User
	userProfileRepo       userRepository.UserProfile
}

func NewOAuthCallback(
	google oauth.Google,
	s3 s3.Client,
	writer *database.Writer,
	authRepo authRepository.Authentication,
	emailVerificationRepo authRepository.EmailVerification,
	sessionRepo authRepository.Session,
	userRepo userRepository.User,
	userProfileRepo userRepository.UserProfile,
) OAuthCallback {
	return &oauthCallback{
		google:                google,
		s3:                    s3,
		writer:                writer,
		authRepo:              authRepo,
		emailVerificationRepo: emailVerificationRepo,
		sessionRepo:           sessionRepo,
		userRepo:              userRepo,
		userProfileRepo:       userProfileRepo,
	}
}

func (uc *oauthCallback) Do(ctx context.Context, input OAuthCallbackInput) (OAuthCallbackOutput, error) {
	var payload *oauth.Payload
	switch input.Provider {
	case oauthModel.OAuthProviderGoogle:
		var err error
		payload, err = uc.google.Validate(ctx, input.Code)
		if err != nil {
			return OAuthCallbackOutput{}, fmt.Errorf("failed to validate google code: %w", err)
		}
	default:
		return OAuthCallbackOutput{}, fmt.Errorf("unsupported provider: %s", input.Provider)
	}

	if payload == nil {
		return OAuthCallbackOutput{}, fmt.Errorf("failed to validate code: %s", input.Code)
	}

	var pictureReader io.ReadSeeker
	if payload.Picture != "" {
		picture, err := http.Get(payload.Picture)
		if err != nil {
			return OAuthCallbackOutput{}, fmt.Errorf("failed to get picture from s3: %w", err)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(picture.Body)

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, picture.Body); err != nil {
			return OAuthCallbackOutput{}, fmt.Errorf("failed to read picture body: %w", err)
		}
		pictureReader = bytes.NewReader(buf.Bytes())
	}

	var session authModel.Session
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		verification := authModel.EmailVerification{
			IntentType: authModel.EmailVerificationIntentTypeSignUp,
			Email:      payload.Email,
		}
		if err := verification.Request(time.Second); err != nil {
			return fmt.Errorf("failed to request email verification: %w", err)
		}
		if err := verification.Verify(); err != nil {
			return fmt.Errorf("failed to verify email verification: %w", err)
		}

		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Create(ctx, &verification); err != nil {
			return fmt.Errorf("failed to create email verification: %w", err)
		}

		user := userModel.User{
			Email: payload.Email,
		}
		if err := uc.userRepo.WithTx(tx.WriterDB()).Create(ctx, &user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		var image commonModel.S3Object
		if pictureReader != nil {
			bucket := config.AWS().S3PublicBucket
			key := fmt.Sprintf("profile_images/%s", user.ID)
			if err := uc.s3.PutObject(ctx, bucket, key, pictureReader); err != nil {
				return fmt.Errorf("failed to upload picture to s3: %w", err)
			}
			image = commonModel.S3Object{
				Bucket: bucket,
				Key:    key,
			}
		}

		profile := userModel.UserProfile{
			UserID: user.ID,
			Name:   payload.Name,
		}
		if !image.IsZero() {
			profile.Image = &image
		}
		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Create(ctx, &profile); err != nil {
			return fmt.Errorf("failed to create user profile: %w", err)
		}

		auth, err := authModel.NewAuthenticationOAuth(*payload)
		if err != nil {
			return fmt.Errorf("failed to create authentication: %w", err)
		}
		auth.UserID = user.ID

		if err := uc.authRepo.WithTx(tx.WriterDB()).Create(ctx, &auth); err != nil {
			return fmt.Errorf("failed to create authentication: %w", err)
		}

		session, err = authModel.NewSession(user.ID)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		return nil
	}); err != nil {
		return OAuthCallbackOutput{}, err
	}

	return OAuthCallbackOutput{Session: session}, nil
}
