package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/message/v1/messagev1connect"
	messagev1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/message/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/message/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

type Message struct {
	send usecase.SendMessage
}

func NewMessage(
	send usecase.SendMessage,
) *Message {
	return &Message{
		send: send,
	}
}

func (h *Message) SendMessage(
	ctx context.Context,
	req *connect.Request[messagev1.SendMessageRequest],
) (*connect.Response[messagev1.SendMessageResponse], error) {
	_, err := h.send.Do(ctx, usecase.SendMessageInput{
		SenderName:   req.Msg.SenderName,
		ReceiverSlug: req.Msg.ReceiverSlug,
		Content:      req.Msg.Content,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&messagev1.SendMessageResponse{})
	return res, nil
}

var _ messagev1connect.MessageServiceHandler = (*Message)(nil)
