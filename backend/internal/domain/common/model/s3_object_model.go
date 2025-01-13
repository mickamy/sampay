package model

import (
	"fmt"

	"gorm.io/gorm"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/lib/ulid"
)

type S3Object struct {
	ID     string
	Bucket string
	Key    string
}

func (m S3Object) URL() string {
	scheme := "https"
	if config.Common().Env == config.Development {
		scheme = "http"
	}
	domain := config.AWS().CloudFrontDomain
	return fmt.Sprintf("%s://%s/%s", scheme, domain, m.Key)
}

func (m S3Object) IsZero() bool {
	return m == S3Object{}
}

func (m *S3Object) BeforeCreate(db *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}
