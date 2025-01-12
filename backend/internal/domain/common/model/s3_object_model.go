package model

import (
	"fmt"

	"mickamy.com/sampay/config"
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
