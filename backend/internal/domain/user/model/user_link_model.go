package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/lib/ulid"
)

var (
	ErrUserLinkInvalidURI = errors.New("invalid uri")
)

type UserLink struct {
	ID           string
	UserID       string
	ProviderType UserLinkProviderType
	URI          string
	QRCodeID     *string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	DisplayAttribute UserLinkDisplayAttribute
	QRCode           *commonModel.S3Object
}

func (m *UserLink) SetQRCode(qrCode *commonModel.S3Object) {
	if qrCode == nil {
		m.QRCodeID = nil
		m.QRCode = nil
		return
	}

	m.QRCode = qrCode
	m.QRCodeID = ptr.Of(qrCode.ID)
}

func (m UserLink) Validate() error {
	if !m.ProviderType.MatchString(m.URI) {
		return errors.Join(ErrUserLinkInvalidURI, fmt.Errorf("uri=[%s]", m.URI))
	}
	return nil
}

func (m *UserLink) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}
