package server

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/database"
	operationv1 "github.com/submaline/services/gen/protocol/operation/v1"
	supervisorv1 "github.com/submaline/services/gen/protocol/supervisor/v1"
	"github.com/submaline/services/gen/protocol/supervisor/v1/supervisorv1connect"
	typesv1 "github.com/submaline/services/gen/protocol/types/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"strconv"
)

var (
	OperationServiceName = zap.String("service", "Operation")
)

type OperationServer struct {
	DB   *database.DBClient // for mariadb
	Auth *auth.Client       // for firebase auth
	//Id     *snowflake.Node    // for id generate
	Rb     *amqp.Connection // for rabbitmq
	Logger *zap.Logger      // for logging

	SvClient *supervisorv1connect.SupervisorServiceClient
}

func (s *OperationServer) FetchOperations(_ context.Context,
	req *connect.Request[operationv1.FetchOperationsRequest],
	stream *connect.ServerStream[operationv1.FetchOperationsResponse]) error {
	funcName := zap.String("func", "FetchOperations")
	logging.LogGrpcFuncCall(s.Logger, OperationServiceName, funcName)
	requesterUserId := req.Header().Get("X-Submaline-UserId")

	// sv用のトークン生成
	adminToken, err := util.GenerateAdminToken()
	if err != nil {
		logging.LogError(s.Logger, OperationServiceName, funcName, "sv用のトークンの生成に失敗しました", err)
		return connect.NewError(connect.CodeUnknown, err)
	}

	recordReq := connect.NewRequest(&supervisorv1.RecordOperationRequest{
		Operations: []*typesv1.Operation{
			{
				//Id:          0,
				Type: typesv1.OperationType_OPERATION_TYPE_FETCH_OPERATIONS,
				//Source:      requesterUserId,
				Destination: []string{requesterUserId},
				//Param1:      "",
				//Param2:      "",
				//Param3:      "",
				// CreatedAt:
			},
		},
	})
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

	_, err = (*s.SvClient).RecordOperation(context.Background(), recordReq)
	if err != nil {
		logging.LogError(s.Logger, OperationServiceName, funcName, "SVにopの配信を依頼できませんでした。", err)
		return connect.NewError(connect.CodeUnknown, err)
	}

	ch, err := s.Rb.Channel()
	if err != nil {
		logging.LogError(s.Logger, OperationServiceName, funcName, "チャンネル生成に失敗しました", err)
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
		logging.LogError(s.Logger, OperationServiceName, funcName, "キューの宣言に失敗しました", err)
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
		logging.LogError(s.Logger, OperationServiceName, funcName, "メッセージの消費に失敗しました", err)
		return connect.NewError(connect.CodeUnknown, err)
	}

	for msg := range messages {
		opId, err := strconv.ParseInt(string(msg.Body), 10, 64)
		if err != nil {
			logging.LogError(s.Logger, OperationServiceName, funcName, "opIdの変換に失敗しました", err)
			return connect.NewError(connect.CodeInternal, err)
		}

		op, err := s.DB.GetOperationWithOperationId(opId)
		if err != nil {
			logging.LogError(s.Logger, OperationServiceName, funcName, "dbからoperationを取得できませんでした", err)
			return connect.NewError(connect.CodeUnknown, err)
		}

		var opMsg *typesv1.Message
		if op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE ||
			op.Type == typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE_RECV {
			m_, err := s.DB.GetMessageWithMessageId(op.Param1)
			if err != nil {
				logging.LogError(s.Logger, OperationServiceName, funcName, "opIdに紐づいているmessageの取得に失敗しました", err)
				return connect.NewError(connect.CodeUnknown, err)
			}

			opMsg = m_
		}

		err = stream.Send(&operationv1.FetchOperationsResponse{
			Operation: op,
			Message:   opMsg,
		})
		if err != nil {
			logging.LogError(s.Logger, OperationServiceName, funcName, "opの配信に失敗しました", err)
		}
		logging.LogInfo(
			s.Logger,
			OperationServiceName,
			funcName,
			fmt.Sprintf("%vにopId: %vを送信しました\ntype: %v\n", requesterUserId, opId, op.Type.String()))
	}

	logging.LogGrpcFuncFinish(s.Logger, OperationServiceName, funcName)
	return nil
}
