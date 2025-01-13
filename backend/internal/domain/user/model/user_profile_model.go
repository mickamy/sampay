package model

import (
	"time"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/ptr"
)

type UserProfile struct {
	UserID    string `gorm:"primaryKey"`
	Name      string
	Bio       *string
	ImageID   *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Image *commonModel.S3Object
}

func (m *UserProfile) SetImage(s3Object *commonModel.S3Object) {
	if s3Object == nil || s3Object.IsZero() {
		m.Image = nil
		m.ImageID = nil
		return
	}

	m.Image = s3Object
	m.ImageID = ptr.Of(s3Object.ID)
}

func (m UserProfile) IsZero() bool {
	return m == UserProfile{}
}
