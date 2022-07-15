package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/db"
	operationv1 "github.com/submaline/services/gen/operation/v1"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	typesv1 "github.com/submaline/services/gen/types/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"strings"
)

type OperationServer struct {
	DB   *db.DBClient // for mariadb
	Auth *auth.Client // for firebase auth
	//Id     *snowflake.Node    // for id generate
	Rb     *amqp.Connection // for rabbitmq
	Logger *zap.Logger      // for logging

	SvClient *supervisorv1connect.SupervisorServiceClient
}

func (s *OperationServer) FetchOperations(ctx context.Context,
	req *connect.Request[operationv1.FetchOperationsRequest],
	stream *connect.ServerStream[operationv1.FetchOperationsResponse]) error {
	//funcName := zap.String("func", "FetchOperations")
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	// sv用のトークン生成
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
		return connect.NewError(connect.CodeUnknown, err)
	}

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{
		Operations: []*typesv1.Operation{
			{
				//Id:          0,
				Type:        typesv1.OperationType_OPERATION_TYPE_FETCH_OPERATIONS,
				Source:      requesterUserId,
				Destination: []string{requesterUserId},
				//Param1:      "",
				//Param2:      "",
				//Param3:      "",
				// CreatedAt:
			},
		},
	})
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminToken.IdToken))

	_, err = (*s.SvClient).RecordOperation(context.Background(), recordReq)
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
		return connect.NewError(connect.CodeUnknown, err)
	}

	ch, err := s.Rb.Channel()
	if err != nil {
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"RabbitMQ チャンネル生成に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return connect.NewError(connect.CodeUnknown, err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		requesterUserId,
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
			"RabbitMQ キュー宣言に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return connect.NewError(connect.CodeUnknown, err)
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
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
			"RabbitMQ メッセージの購読に失敗しました",
			nil,
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return connect.NewError(connect.CodeUnknown, err)
	}

	for {
		select {
		case <-ctx.Done():
			logging.Info(s.Logger, req.Spec().Procedure, "userとの接続が切れました。")
			return nil
		case msg := <-messages:
			opId, err := strconv.ParseInt(string(msg.Body), 10, 64)
			if err != nil {
				// log
				if e_ := logging.ErrD(
					s.Logger,
					req.Spec().Procedure,
					err,
					"Operation.idの変換に失敗しました",
					nil,
					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
					log.Println(e_)
				}
				return connect.NewError(connect.CodeInternal, err)
			}

			op, err := s.DB.GetOperationWithOperationId(opId)
			if err != nil {
				// log
				if e_ := logging.ErrD(
					s.Logger,
					req.Spec().Procedure,
					err,
					"Operationのデータ取得に失敗しました",
					nil,
					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
					log.Println(e_)
				}
				return connect.NewError(connect.CodeUnknown, err)
			}

			var opMsg *typesv1.Message
			if op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE ||
				op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE_RECV {
				m_, err := s.DB.GetMessageWithMessageId(op.Param1)
				if err != nil {
					// log
					if e_ := logging.ErrD(
						s.Logger,
						req.Spec().Procedure,
						err,
						"Operationに付属するMessageの取得に失敗しました",
						nil,
						os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
						log.Println(e_)
					}
					return connect.NewError(connect.CodeUnknown, err)
				}

				opMsg = m_
			}

			err = stream.Send(&operationv1.FetchOperationsResponse{
				Operation: op,
				Message:   opMsg,
			})
			if err != nil {
				// log
				if e_ := logging.ErrD(
					s.Logger,
					req.Spec().Procedure,
					err,
					"Operationの配信に失敗しました",
					[]logging.DiscordRichMessageEmbedField{
						logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", op.Id), false),
						logging.GenerateDiscordRichMsgField("type", fmt.Sprintf("%v", op.Type.String()), false),
						logging.GenerateDiscordRichMsgField("source", fmt.Sprintf("%v", op.Source), false),
						logging.GenerateDiscordRichMsgField("dest", fmt.Sprintf("%v", strings.Join(op.Destination, ", ")), false),
					},
					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
					log.Println(e_)
				}
			} else { // log
				if e_ := logging.InfoD(
					s.Logger,
					req.Spec().Procedure,
					"Operationを配信しました",
					[]logging.DiscordRichMessageEmbedField{
						logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", op.Id), false),
					},
					os.Getenv("DISCORD_WEBHOOK_URL"),
				); e_ != nil {
					log.Println(e_)
				}
			}
		}
	}

	//select {
	//case <-ctx.Done():
	//	return nil
	//default:
	//
	//	for msg := range messages {
	//		opId, err := strconv.ParseInt(string(msg.Body), 10, 64)
	//		if err != nil {
	//			// log
	//			if e_ := logging.ErrD(
	//				s.Logger,
	//				req.Spec().Procedure,
	//				err,
	//				"Operation.idの変換に失敗しました",
	//				nil,
	//				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
	//				log.Println(e_)
	//			}
	//			return connect.NewError(connect.CodeInternal, err)
	//		}
	//
	//		op, err := s.DB.GetOperationWithOperationId(opId)
	//		if err != nil {
	//			// log
	//			if e_ := logging.ErrD(
	//				s.Logger,
	//				req.Spec().Procedure,
	//				err,
	//				"Operationのデータ取得に失敗しました",
	//				nil,
	//				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
	//				log.Println(e_)
	//			}
	//			return connect.NewError(connect.CodeUnknown, err)
	//		}
	//
	//		var opMsg *typesv1.Message
	//		if op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE ||
	//			op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE_RECV {
	//			m_, err := s.DB.GetMessageWithMessageId(op.Param1)
	//			if err != nil {
	//				// log
	//				if e_ := logging.ErrD(
	//					s.Logger,
	//					req.Spec().Procedure,
	//					err,
	//					"Operationに付属するMessageの取得に失敗しました",
	//					nil,
	//					os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
	//					log.Println(e_)
	//				}
	//				return connect.NewError(connect.CodeUnknown, err)
	//			}
	//
	//			opMsg = m_
	//		}
	//
	//		err = stream.Send(&operationv1.FetchOperationsResponse{
	//			Operation: op,
	//			Message:   opMsg,
	//		})
	//		if err != nil {
	//			// log
	//			if e_ := logging.ErrD(
	//				s.Logger,
	//				req.Spec().Procedure,
	//				err,
	//				"Operationの配信に失敗しました",
	//				[]logging.DiscordRichMessageEmbedField{
	//					logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", op.Id), false),
	//					logging.GenerateDiscordRichMsgField("type", fmt.Sprintf("%v", op.Type.String()), false),
	//					logging.GenerateDiscordRichMsgField("source", fmt.Sprintf("%v", op.Source), false),
	//					logging.GenerateDiscordRichMsgField("dest", fmt.Sprintf("%v", strings.Join(op.Destination, ", ")), false),
	//				},
	//				os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
	//				log.Println(e_)
	//			}
	//		} else { // log
	//			if e_ := logging.InfoD(
	//				s.Logger,
	//				req.Spec().Procedure,
	//				"Operationを配信しました",
	//				[]logging.DiscordRichMessageEmbedField{
	//					logging.GenerateDiscordRichMsgField("opId", fmt.Sprintf("%v", op.Id), false),
	//				},
	//				os.Getenv("DISCORD_WEBHOOK_URL"),
	//			); e_ != nil {
	//				log.Println(e_)
	//			}
	//		}
	//	}
	//}

	//return nil
}
