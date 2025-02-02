package response

import (
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"

	"mickamy.com/sampay/internal/lib/paging"
)

func NewNextPage(page paging.NextPage) *commonv1.NextPage {
	return &commonv1.NextPage{
		Index: int32(page.Index),
		Limit: int32(page.Limit),
		Total: int32(page.Total),
	}
}
