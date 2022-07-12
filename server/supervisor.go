package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/bwmarrin/snowflake"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/database"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"time"
)

const (
	defaultDisplayName = ""
	defaultIconPath    = ""
)

var (
	SupervisorServiceName = zap.String("service", "Supervisor")
)

type SupervisorServer struct {
	DB     *database.DBClient // for mariadb
	Auth   *auth.Client       // for firebase auth
	Id     *snowflake.Node    // for id generate
	Rb     *amqp.Connection   // for rabbitmq
	Logger *zap.Logger        // for logging
}

func (s *SupervisorServer) CreateAccount(_ context.Context,
	req *connect.Request[supervisorv1.CreateAccountRequest]) (
	*connect.Response[supervisorv1.CreateAccountResponse], error) {
	funcName := zap.String("func", "CreateAccount")
	logging.LogGrpcFuncCall(s.Logger, SupervisorServiceName, funcName)

	// firebaseのトークンにadminクレームが入っていれば、その情報がインターセプターで挿入されてるはず
	if !util.ParseBool(req.Header().Get("X-Submaline-Admin")) {
		err := ErrAdminOnly
		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)

		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	// firebaseにユーザーが存在するかを確認
	user, err := s.Auth.GetUser(context.Background(), req.Msg.Account.UserId)
	if err != nil {
		if auth.IsUserNotFound(err) {
			logging.LogError(
				s.Logger,
				SupervisorServiceName,
				funcName,
				"ユーザーが存在しません",
				err)

			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"firebaseで不明なエラー",
			err)

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// databaseに作成
	account, err := s.DB.CreateAccount(user.UID, user.Email)
	if err != nil {
		// todo : already exists
		// Error 1062: Duplicate entry

		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"データベースでエラー",
			err)

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&supervisorv1.CreateAccountResponse{Account: account})

	logging.LogGrpcFuncFinish(s.Logger, SupervisorServiceName, funcName)

	return res, nil
}

func (s *SupervisorServer) CreateProfile(_ context.Context,
	req *connect.Request[supervisorv1.CreateProfileRequest]) (
	*connect.Response[supervisorv1.CreateProfileResponse], error) {
	funcName := zap.String("func", "CreateProfile")

	logging.LogGrpcFuncCall(s.Logger, SupervisorServiceName, funcName)

	// 権限確認
	if !util.ParseBool(req.Header().Get("X-Peg-Admin")) {
		err := ErrAdminOnly

		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)

		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("admin only"))
	}

	// firebaseにユーザーが存在するか
	user, err := s.Auth.GetUser(context.Background(), req.Msg.Profile.UserId)
	if err != nil {
		if auth.IsUserNotFound(err) {

			logging.LogError(
				s.Logger,
				SupervisorServiceName,
				funcName,
				"ユーザー存在しない",
				err)

			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// databaseに作成
	profile, err := s.DB.CreateProfile(user.UID, defaultDisplayName, defaultIconPath)
	if err != nil {
		// todo : already exists
		// Error 1062: Duplicate entry
		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&supervisorv1.CreateProfileResponse{Profile: profile})

	logging.LogGrpcFuncFinish(s.Logger, SupervisorServiceName, funcName)

	return res, nil
}

func (s *SupervisorServer) RecordOperation(_ context.Context,
	req *connect.Request[supervisorv1.RecordOperationRequest]) (
	*connect.Response[supervisorv1.RecordOperationResponse], error) {
	funcName := zap.String("func", "RecordOperation")
	logging.LogGrpcFuncCall(s.Logger, SupervisorServiceName, funcName)

	// 権限確認
	if !util.ParseBool(req.Header().Get("X-Peg-Admin")) {
		err := ErrAdminOnly
		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("admin only"))
	}

	ch, err := s.Rb.Channel()
	if err != nil {
		logging.LogError(
			s.Logger,
			SupervisorServiceName,
			funcName,
			"",
			err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	defer ch.Close()

	// supervisorは、お願いされたものを記録するだけ。
	for _, op := range req.Msg.Operations {
		opId := s.Id.Generate().Int64()
		// op本体の記録

		// 時間強制書き換え
		_, err := s.DB.CreateOperation(opId, op.Type, op.Source, op.Param1, op.Param2, op.Param3, time.Now())
		if err != nil {
			logging.LogError(
				s.Logger,
				SupervisorServiceName,
				funcName,
				"",
				err)
			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		// destinationを記録する前にopを記録しているので、宛先が出鱈目でもopだけは保存されてしまう。
		// destinationのop_idはfkなので、先に入れることはできない。
		// 先にdestをチェックするか、、、

		// 宛先の記録
		err = s.DB.CreateOperationDestination(opId, op.Destination)
		if err != nil {
			logging.LogError(
				s.Logger,
				SupervisorServiceName,
				funcName,
				"",
				err)
			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		for _, dest := range op.Destination {
			q, err := ch.QueueDeclare(
				dest,
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				logging.LogError(
					s.Logger,
					SupervisorServiceName,
					funcName,
					"",
					err)
				return nil, connect.NewError(connect.CodeUnknown, err)
			}

			err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("%v", opId)),
				})
			if err != nil {
				logging.LogError(
					s.Logger,
					SupervisorServiceName,
					funcName,
					"",
					err)
				return nil, connect.NewError(connect.CodeUnknown, err)
			}

			logging.LogInfo(
				s.Logger,
				SupervisorServiceName,
				funcName,
				fmt.Sprintf("MQに%v宛のop: %vの到着をお知らせしました", dest, opId))

		}
	}

	res := connect.NewResponse(&supervisorv1.RecordOperationResponse{})

	return res, nil
}
