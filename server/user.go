package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/submaline/services/db"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	userv1 "github.com/submaline/services/gen/user/v1"
	"github.com/submaline/services/logging"
	"go.uber.org/zap"
)

var (
	UserServiceName = zap.String("serviceName", "User")
)

type UserServer struct {
	DB       *db.DBClient
	Auth     *auth.Client
	Logger   *zap.Logger
	SvClient *supervisorv1connect.SupervisorServiceClient
}

func (s *UserServer) GetAccount(_ context.Context,
	req *connect.Request[userv1.GetAccountRequest]) (
	*connect.Response[userv1.GetAccountResponse], error) {
	funcName := zap.String("funcName", "GetAccount")
	logging.LogGrpcFuncCall(s.Logger, UserServiceName, funcName)
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	account, err := s.DB.GetAccount(requesterUserId)
	if err != nil {
		logging.LogError(s.Logger, UserServiceName, funcName, "failed to get account data", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	logging.LogGrpcFuncFinish(s.Logger, UserServiceName, funcName)
	return connect.NewResponse(&userv1.GetAccountResponse{
		Account: account,
	}), nil
}

func (s *UserServer) UpdateAccount(_ context.Context,
	_ *connect.Request[userv1.UpdateAccountRequest]) (
	*connect.Response[userv1.UpdateAccountResponse], error) {
	funcName := zap.String("funcName", "UpdateAccount")
	logging.LogGrpcFuncCall(s.Logger, UserServiceName, funcName)
	//requesterUserId := req.Header().Get("X-Submaline-UserId")

	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented: UpdateAccount"))
}

func (s *UserServer) GetProfile(_ context.Context,
	req *connect.Request[userv1.GetProfileRequest]) (
	*connect.Response[userv1.GetProfileResponse], error) {
	funcName := zap.String("funcName", "GetProfile")
	logging.LogGrpcFuncCall(s.Logger, UserServiceName, funcName)

	prof, err := s.DB.GetProfile(req.Msg.UserId)
	if err != nil {
		logging.LogError(s.Logger, UserServiceName, funcName, "", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	logging.LogGrpcFuncFinish(s.Logger, UserServiceName, funcName)
	return connect.NewResponse(&userv1.GetProfileResponse{
		Profile: prof,
	}), nil
}

func (s *UserServer) UpdateProfile(_ context.Context,
	req *connect.Request[userv1.UpdateProfileRequest]) (
	*connect.Response[userv1.UpdateProfileResponse], error) {
	funcName := zap.String("funName", "UpdateProfile")
	logging.LogGrpcFuncCall(s.Logger, UserServiceName, funcName)
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	prof, err := s.DB.UpdateProfile(requesterUserId, req.Msg)
	if err != nil {
		logging.LogError(s.Logger, UserServiceName, funcName, "", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	logging.LogGrpcFuncFinish(s.Logger, UserServiceName, funcName)
	return connect.NewResponse(&userv1.UpdateProfileResponse{Profile: prof}), nil
}
