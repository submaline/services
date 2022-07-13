package interceptor

import (
	"context"
	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type logInterceptor struct {
	*zap.Logger
}

func (i *logInterceptor) WrapStreamContext(ctx context.Context) context.Context {
	return ctx
}

func NewLogInterceptor(l *zap.Logger) *logInterceptor {
	return &logInterceptor{l}
}

// WrapUnary client 1 -> server 1
func (i *logInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		// メソッドが呼び出されたことを通知
		i.Logger.Info(
			"UnaryMethodCall",
			zap.String("method", request.Spec().Procedure),
		)

		res, err := next(ctx, request)
		if err != nil {
			return nil, err
		}

		i.Logger.Info(
			"UnaryMethodFinish",
			zap.String("method", request.Spec().Procedure),
		)
		return res, nil
	}
}

// WrapStreamSender ?
func (i *logInterceptor) WrapStreamSender(_ context.Context, s connect.Sender) connect.Sender {
	return &logStreamSender{
		Sender: s,
		Logger: i.Logger,
	}
}

type logStreamSender struct {
	connect.Sender
	*zap.Logger
}

func (s *logStreamSender) Send(msg any) error {
	if err := s.Sender.Send(msg); err != nil {
		return err
	}

	s.Logger.Info(
		"StreamMethodCall",
		zap.String("stream", "send"),
		zap.String("method", s.Spec().Procedure),
	)

	return nil
}

// WrapStreamReceiver client 1 -> server stream
func (i *logInterceptor) WrapStreamReceiver(_ context.Context, r connect.Receiver) connect.Receiver {
	return &logStreamReceiver{
		Receiver: r,
		Logger:   i.Logger,
	}
}

type logStreamReceiver struct {
	connect.Receiver
	*zap.Logger
}

// Receive Streamのリクエストを受信した呼ばれると思われる
func (r *logStreamReceiver) Receive(msg any) error {
	if err := r.Receiver.Receive(msg); err != nil {
		return err
	}

	r.Logger.Info(
		"StreamMethodCall",
		zap.String("stream", "receive"),
		zap.String("method", r.Spec().Procedure),
	)

	return nil
}
