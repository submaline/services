package interceptors

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"strings"
)

type authInterceptor struct {
	authClient *auth.Client
}

func NewAuthInterceptor(fbAuthClient *auth.Client) *authInterceptor {
	return &authInterceptor{authClient: fbAuthClient}
}

func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		// インターセプター（前処理）
		// "Bearer e..."
		idTokenRaw := request.Header().Get("Authorization")
		if idTokenRaw == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("need Bearer token"))
		}
		// "e..."
		idToken := strings.Replace(idTokenRaw, "Bearer ", "", 1)
		// 検証
		token, err := i.authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		// カスタムクレーム検証
		claims := token.Claims
		// 管理者チェック
		if admin, ok := claims["admin"]; ok {
			if admin.(bool) {
				request.Header().Set("X-Submaline-Admin", "true")
			}
		} else {
			request.Header().Set("X-Submaline-Admin", "false")
		}

		// 本来の処理で使用するのでデータをくっつけてあげる。
		request.Header().Set("X-Submaline-UserId", token.UID)

		// 本来の処理
		res, err := next(ctx, request)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
}

func (*authInterceptor) WrapStreamContext(ctx context.Context) context.Context {
	return ctx
}

func (i *authInterceptor) WrapStreamSender(_ context.Context, s connect.Sender) connect.Sender {
	return &authSender{
		Sender: s,
	}
}

func (i *authInterceptor) WrapStreamReceiver(_ context.Context, r connect.Receiver) connect.Receiver {
	return &authReceiver{
		Receiver:   r,
		AuthClient: i.authClient,
	}
}

type authSender struct {
	connect.Sender
}

func (s *authSender) Send(msg any) error {
	return s.Sender.Send(msg)
}

type authReceiver struct {
	connect.Receiver
	AuthClient *auth.Client
}

func (r *authReceiver) Receive(msg any) error {
	if err := r.Receiver.Receive(msg); err != nil {
		return err
	}

	// インターセプター（前処理）
	// "Bearer e..."
	idTokenRaw := r.Header().Get("Authorization")
	if idTokenRaw == "" {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("need Bearer token"))
	}
	// "e..."
	idToken := strings.Replace(idTokenRaw, "Bearer ", "", 1)
	// 検証
	token, err := r.AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	r.Header().Set("X-Submaline-UserId", token.UID)

	//r.logger.Printf("<Receive>%s: req: %v", r.Spec().Procedure, msg)
	//log.Printf("<Receive>%s", r.Header().Get("Authorization"))

	return nil
}
