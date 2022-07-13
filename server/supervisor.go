package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/bwmarrin/snowflake"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/db"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

const (
	defaultDisplayName = "unknown"
	defaultIconPath    = "d2f60d27-a8b1-4631-8cc8-4f3a24155599"
)

type SupervisorServer struct {
	DB     *db.DBClient     // for mariadb
	Auth   *auth.Client     // for firebase auth
	Id     *snowflake.Node  // for id generate
	Rb     *amqp.Connection // for rabbitmq
	Logger *zap.Logger      // for logging
}

func (s *SupervisorServer) CreateAccount(_ context.Context,
	req *connect.Request[supervisorv1.CreateAccountRequest]) (
	*connect.Response[supervisorv1.CreateAccountResponse], error) {

	// firebaseのトークンにadminクレームが入っていれば、その情報がインターセプターで挿入されてるはず
	if !util.ParseBool(req.Header().Get("X-Submaline-Admin")) {
		err := ErrAdminOnly
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者以外がアカウント作成を試みました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	// firebaseにユーザーが存在するかを確認
	user, err := s.Auth.GetUser(context.Background(), req.Msg.Account.UserId)
	if err != nil {
		if auth.IsUserNotFound(err) {
			// log
			if e_ := logging.ErrD(
				s.Logger,
				req.Spec().Procedure,
				err,
				"Firebaseに登録されていないユーザーのアカウントの作成はできません",
				[]logging.DiscordRichMessageEmbedField{
					logging.GenerateDiscordRichMsgField("作成しようとしたUser", req.Msg.Account.UserId, false),
				},
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(e_)
			}
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			ErrMsgFailedToGetUserDatFromFirebase,
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// databaseに作成
	account, err := s.DB.CreateAccount(user.UID, user.Email)
	if err != nil {
		// todo : already exists
		// Error 1062: Duplicate entry

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"アカウントを作成できませんでした",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&supervisorv1.CreateAccountResponse{Account: account})

	return res, nil
}

func (s *SupervisorServer) CreateProfile(_ context.Context,
	req *connect.Request[supervisorv1.CreateProfileRequest]) (
	*connect.Response[supervisorv1.CreateProfileResponse], error) {

	// 権限確認
	if !util.ParseBool(req.Header().Get("X-Submaline-Admin")) {
		err := ErrAdminOnly

		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者以外がプロフィール作成を試みました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	// firebaseにユーザーが存在するか
	user, err := s.Auth.GetUser(context.Background(), req.Msg.Profile.UserId)
	if err != nil {
		if auth.IsUserNotFound(err) {

			// log
			if e_ := logging.ErrD(
				s.Logger,
				req.Spec().Procedure,
				err,
				"Firebaseに登録されていないユーザーのプロフィールは作成できません",
				[]logging.DiscordRichMessageEmbedField{
					logging.GenerateDiscordRichMsgField("作成しようとしたUser", req.Msg.Profile.UserId, false),
				},
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(e_)
			}

			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			ErrMsgFailedToGetUserDatFromFirebase,
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// databaseに作成
	profile, err := s.DB.CreateProfile(user.UID, defaultDisplayName, defaultIconPath)
	if err != nil {
		// todo : already exists
		// Error 1062: Duplicate entry

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"プロフィールの作成に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&supervisorv1.CreateProfileResponse{Profile: profile})

	return res, nil
}

func (s *SupervisorServer) RecordOperation(_ context.Context,
	req *connect.Request[supervisorv1.RecordOperationRequest]) (
	*connect.Response[supervisorv1.RecordOperationResponse], error) {

	// 権限確認
	if !util.ParseBool(req.Header().Get("X-Submaline-Admin")) {
		err := ErrAdminOnly

		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"管理者以外がOperation記録を試みました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}

		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	ch, err := s.Rb.Channel()
	if err != nil {
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			ErrMsgFailedToGetUserDatFromFirebase,
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
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
			// log
			if e_ := logging.ErrD(
				s.Logger,
				req.Spec().Procedure,
				err,
				"データベースにOperationを書き込めませんでした",
				[]logging.DiscordRichMessageEmbedField{
					logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", opId), false),
				},
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(e_)
			}
			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		// destinationを記録する前にopを記録しているので、宛先が出鱈目でもopだけは保存されてしまう。
		// destinationのop_idはfkなので、先に入れることはできない。
		// 先にdestをチェックするか、、、

		// 宛先の記録
		err = s.DB.CreateOperationDestination(opId, op.Destination)
		if err != nil {
			// log
			if e_ := logging.ErrD(
				s.Logger,
				req.Spec().Procedure,
				err,
				"Operationの宛先記録に失敗しました",
				[]logging.DiscordRichMessageEmbedField{
					logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", opId), false),
				},
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(e_)
			}
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
				// log
				if e_ := logging.ErrD(
					s.Logger,
					req.Spec().Procedure,
					err,
					"RabbitMQ キューの宣言に失敗しました",
					nil,
					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
					log.Println(e_)
				}
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
				// log
				if e_ := logging.ErrD(
					s.Logger,
					req.Spec().Procedure,
					err,
					"RabbitMQ メッセージの送信に失敗しました",
					nil,
					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
					log.Println(e_)
				}
				return nil, connect.NewError(connect.CodeUnknown, err)
			}

			// log
			if e_ := logging.InfoD(
				s.Logger,
				req.Spec().Procedure,
				"Operationを正常に記録しました",
				[]logging.DiscordRichMessageEmbedField{
					logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", opId), false),
					logging.GenerateDiscordRichMsgField("opType", op.Type.String(), false),
				},
				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
				log.Println(err)
			}

		}
	}

	res := connect.NewResponse(&supervisorv1.RecordOperationResponse{})

	return res, nil
}
