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
	"strconv"
)

var (
	AuthServiceName = zap.String("service", "Auth")
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
	funcName := zap.String("func", "LoginWithEmail")
	logging.LogGrpcFuncCall(s.Logger, AuthServiceName, funcName)

	token, err := util.GenerateToken(req.Msg.Email, req.Msg.Password)
	if err != nil {
		logging.LogError(s.Logger, AuthServiceName, funcName, "", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	expiresIn, err := strconv.ParseInt(token.ExpiresIn, 10, 64)
	if err != nil {
		logging.LogError(s.Logger, AuthServiceName, funcName, "failed to parse expiresIn", err)
	}

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{Operations: []*typesv1.Operation{
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_LOGIN_WITH_EMAIL,
			Source:      token.LocalId,
			Destination: []string{token.LocalId},
			Param1:      "",
			Param2:      "",
			Param3:      "",
			CratedAt:    timestamppb.Now(),
		},
	}})

	// sv用のトークン生成
	adminToken, err := util.GenerateToken(os.Getenv("SUBMALINE_ADMIN_FB_EMAIL"), os.Getenv("SUBMALINE_ADMIN_FB_PASSWORD"))
	if err != nil {
		logging.LogError(s.Logger, TalkServiceName, funcName, "sv用のトークンの生成に失敗しました", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminToken.IdToken))
	go func() {
		_, err = (*s.SvClient).RecordOperation(context.Background(), recordReq)
		if err != nil {
			if err != nil {
				log.Println(err)
			}
		}
	}()

	resp := connect.NewResponse(&authv1.LoginWithEmailResponse{
		AuthToken: &typesv1.AuthToken{
			Token:        token.IdToken,
			ExpiresIn:    expiresIn,
			RefreshToken: token.RefreshToken,
		},
	})

	logging.LogGrpcFuncFinish(s.Logger, AuthServiceName, funcName)
	return resp, nil
}

func (s *AuthServer) UpdatePassword(_ context.Context,
	_ *connect.Request[authv1.UpdatePasswordRequest]) (
	*connect.Response[authv1.UpdatePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented: UpdatePassword"))
}
