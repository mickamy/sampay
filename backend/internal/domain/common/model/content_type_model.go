package model

import (
	"fmt"

	"mickamy.com/sampay/internal/lib/either"
)

type ContentType string

func (m ContentType) String() string {
	return string(m)
}

const (
	ContentTypeAudioMPEG ContentType = "audio/mpeg"
	ContentTypeImageJPEG ContentType = "image/jpeg"
	ContentTypeImageJPG  ContentType = "image/jpg"
	ContentTypeImagePNG  ContentType = "image/png"
	ContentTypeImageGIF  ContentType = "image/gif"
	ContentTypeImageBMP  ContentType = "image/bmp"
	ContentTypeTextPlain ContentType = "text/plain"
	ContentTypeTextCSV   ContentType = "text/csv"
	ContentTypeVideoMP4  ContentType = "video/mp4"
	ContentTypeVideoMPEG ContentType = "video/mpeg"
)

func MustNewContentType(s string) ContentType {
	return either.Must(NewContentType(s))
}

func NewContentType(s string) (ContentType, error) {
	for _, messageType := range []ContentType{
		ContentTypeAudioMPEG,
		ContentTypeImageJPEG,
		ContentTypeImageJPG,
		ContentTypeImagePNG,
		ContentTypeImageGIF,
		ContentTypeImageBMP,
		ContentTypeTextPlain,
		ContentTypeTextCSV,
		ContentTypeVideoMP4,
		ContentTypeVideoMPEG,
	} {
		if s == messageType.String() {
			return messageType, nil
		}
	}

	return "", fmt.Errorf("unsupported content type: [%s]", s)
}
