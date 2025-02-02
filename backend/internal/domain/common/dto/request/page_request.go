package request

import (
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"

	"mickamy.com/sampay/internal/lib/paging"
	"mickamy.com/sampay/internal/lib/ptr"
)

func NewPage(pb *commonv1.Page) *paging.Page {
	if pb == nil {
		return nil
	}

	return ptr.Of(paging.New(int(pb.Index), int(pb.Limit)))
}
