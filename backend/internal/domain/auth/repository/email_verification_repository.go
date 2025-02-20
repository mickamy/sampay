package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/auth/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type EmailVerification interface {
	Create(ctx context.Context, m *model.EmailVerification) error
	FindByEmail(ctx context.Context, email string, scope ...database.Scope) (*model.EmailVerification, error)
	FindByRequestedTokenAndPinCode(ctx context.Context, token, pinCode string, scope ...database.Scope) (*model.EmailVerification, error)
	FindByVerifiedToken(ctx context.Context, token string, scope ...database.Scope) (*model.EmailVerification, error)
	Update(ctx context.Context, m *model.EmailVerification) error
	WithTx(tx *database.DB) EmailVerification
}

type emailVerification struct {
	db *database.DB
}

func NewEmailVerification(db *database.DB) EmailVerification {
	return &emailVerification{db: db}
}

func (repo *emailVerification) Create(ctx context.Context, m *model.EmailVerification) error {
	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *emailVerification) FindByEmail(ctx context.Context, email string, scopes ...database.Scope) (*model.EmailVerification, error) {
	var m model.EmailVerification
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		First(&m, "email = ?", email).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *emailVerification) FindByRequestedTokenAndPinCode(ctx context.Context, token, pinCode string, scopes ...database.Scope) (*model.EmailVerification, error) {
	var m model.EmailVerification
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		InnerJoins("Requested").
		First(&m, `"Requested".token = ? AND "Requested".pin_code = ?`, token, pinCode).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *emailVerification) FindByVerifiedToken(ctx context.Context, token string, scopes ...database.Scope) (*model.EmailVerification, error) {
	var m model.EmailVerification
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		InnerJoins("Verified").
		First(&m, `"Verified".token = ?`, token).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *emailVerification) Update(ctx context.Context, m *model.EmailVerification) error {
	return repo.db.WithContext(ctx).Save(m).Error
}

func (repo *emailVerification) WithTx(tx *database.DB) EmailVerification {
	return &emailVerification{db: tx}
}

func EmailVerificationInnerJoinRequested(db *database.DB) *database.DB {
	return &database.DB{DB: db.InnerJoins("Requested")}
}

func EmailVerificationJoinVerified(db *database.DB) *database.DB {
	return &database.DB{DB: db.Joins("Verified")}
}

func EmailVerificationJoinConsumed(db *database.DB) *database.DB {
	return &database.DB{DB: db.Joins("Consumed")}
}

func EmailVerificationNotConsumed(db *database.DB) *database.DB {
	return &database.DB{DB: db.
		Joins("LEFT OUTER JOIN consumed_email_verifications consumed ON email_verifications.id = consumed.email_verification_id").
		Where("consumed.email_verification_id IS NULL"),
	}
}
