package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/submaline/services/db"
	authv1 "github.com/submaline/services/gen/auth/v1"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	typesv1 "github.com/submaline/services/gen/types/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"time"
)

type AuthServer struct {
	DB       *db.DBClient
	Auth     *auth.Client
	Logger   *zap.Logger
	SvClient *supervisorv1connect.SupervisorServiceClient
}

func (s *AuthServer) LoginWithEmail(_ context.Context,
	req *connect.Request[authv1.LoginWithEmailRequest]) (
	*connect.Response[authv1.LoginWithEmailResponse], error) {

	token, err := util.GenerateToken(req.Msg.Email, req.Msg.Password, false)
	if err != nil {
		// todo : firebase invalid ...
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	expiresIn := token.ExpiresAt.Sub(time.Now()).Seconds()

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{Operations: []*typesv1.Operation{
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_LOGIN_WITH_EMAIL,
			Source:      token.UID,
			Destination: []string{token.UID},
			Param1:      "",
			Param2:      "",
			Param3:      "",
			CreatedAt:   timestamppb.Now(),
		},
	}})

	// sv用のトークン生成
	adminToken, err := util.GenerateToken(os.Getenv("SUBMALINE_ADMIN_FB_EMAIL"), os.Getenv("SUBMALINE_ADMIN_FB_PASSWORD"), false)
	if err != nil {
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者トークンの発行に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminToken.IdToken))
	go func() {
		_, err = (*s.SvClient).RecordOperation(context.Background(), recordReq)
		if err != nil {
			if e_ := logging.ErrD(
				s.Logger,
				req.Spec().Procedure,
				err,
				"Operationの記録に失敗しました",
				nil,
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(e_)
			}
		}
	}()

	resp := connect.NewResponse(&authv1.LoginWithEmailResponse{
		AuthToken: &typesv1.AuthToken{
			Token:        token.IdToken,
			ExpiresIn:    int64(expiresIn),
			RefreshToken: token.Refresh,
		},
	})

	return resp, nil
}

func (s *AuthServer) UpdatePassword(_ context.Context,
	_ *connect.Request[authv1.UpdatePasswordRequest]) (
	*connect.Response[authv1.UpdatePasswordResponse], error) {
	// todo : implement
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented: UpdatePassword"))
}
