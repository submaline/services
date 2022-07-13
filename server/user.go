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
	"log"
	"os"
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
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	account, err := s.DB.GetAccount(requesterUserId)
	if err != nil {

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"データベースからアカウントを取得できませんでした",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&userv1.GetAccountResponse{
		Account: account,
	}), nil
}

func (s *UserServer) UpdateAccount(_ context.Context,
	_ *connect.Request[userv1.UpdateAccountRequest]) (
	*connect.Response[userv1.UpdateAccountResponse], error) {

	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented: UpdateAccount"))
}

func (s *UserServer) GetProfile(_ context.Context,
	req *connect.Request[userv1.GetProfileRequest]) (
	*connect.Response[userv1.GetProfileResponse], error) {

	prof, err := s.DB.GetProfile(req.Msg.UserId)
	if err != nil {

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"データベースからプロフィールの取得に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&userv1.GetProfileResponse{
		Profile: prof,
	}), nil
}

func (s *UserServer) UpdateProfile(_ context.Context,
	req *connect.Request[userv1.UpdateProfileRequest]) (
	*connect.Response[userv1.UpdateProfileResponse], error) {
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	prof, err := s.DB.UpdateProfile(requesterUserId, req.Msg)
	if err != nil {

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"データベースでユーザープロフィールの更新に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&userv1.UpdateProfileResponse{Profile: prof}), nil
}
