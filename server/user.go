package server

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/submaline/services/db"
	authv1 "github.com/submaline/services/gen/auth/v1"
	"github.com/submaline/services/gen/auth/v1/authv1connect"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	typesv1 "github.com/submaline/services/gen/types/v1"
	userv1 "github.com/submaline/services/gen/user/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"log"
	"os"
)

type UserServer struct {
	DB     *db.DBClient
	Logger *zap.Logger

	AuthClient *authv1connect.AuthServiceClient
	SvClient   *supervisorv1connect.SupervisorServiceClient
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

	adminToken, err := util.GenerateToken(os.Getenv("SUBMALINE_ADMIN_FB_EMAIL"), os.Getenv("SUBMALINE_ADMIN_FB_PASSWORD"), false)
	if err != nil {

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者トークンの生成に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{Operations: []*typesv1.Operation{
		{
			Id:          0, // svで自動決定
			Type:        typesv1.OperationType_OPERATION_TYPE_GET_ACCOUNT,
			Source:      requesterUserId,
			Destination: []string{requesterUserId},
			Param1:      "",
			Param2:      "",
			Param3:      "",
			CreatedAt:   nil, // svで自動決定
		},
	}})
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %v", adminToken.IdToken))

	go func() {
		_, err = (*s.SvClient).RecordOperation(
			context.Background(),
			recordReq,
		)
		if err != nil {
			// log
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

	requesterUserId := req.Header().Get("X-Submaline-UserId")

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

	// sv用トークン
	adminTokenResp, err := (*s.AuthClient).LoginWithEmail(context.Background(), connect.NewRequest(&authv1.LoginWithEmailRequest{
		Email:    os.Getenv("SUBMALINE_ADMIN_FB_EMAIL"),
		Password: os.Getenv("SUBMALINE_ADMIN_FB_PASSWORD"),
	}))
	if err != nil {
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者トークンの生成に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{Operations: []*typesv1.Operation{
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_GET_PROFILE,
			Source:      requesterUserId,
			Destination: []string{requesterUserId},
		},
	}})
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminTokenResp.Msg.AuthToken.Token))
	go func() {
		_, err = (*s.SvClient).RecordOperation(
			context.Background(),
			recordReq,
		)
		if err != nil {
			// log
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

	adminToken, err := util.GenerateToken(os.Getenv("SUBMALINE_ADMIN_FB_EMAIL"), os.Getenv("SUBMALINE_ADMIN_FB_PASSWORD"), false)
	if err != nil {

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者トークンの生成に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// todo : notify friends
	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{Operations: []*typesv1.Operation{
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_UPDATE_PROFILE,
			Source:      requesterUserId,
			Destination: []string{requesterUserId},
		},
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_UPDATE_PROFILE_RECV,
			Source:      requesterUserId,
			Destination: []string{requesterUserId},
		},
	}})
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminToken.IdToken))
	go func() {
		_, err = (*s.SvClient).RecordOperation(
			context.Background(),
			recordReq,
		)
		if err != nil {
			// log
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

	return connect.NewResponse(&userv1.UpdateProfileResponse{Profile: prof}), nil
}
