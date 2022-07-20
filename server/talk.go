package server

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
	"github.com/submaline/services/db"
	authv1 "github.com/submaline/services/gen/auth/v1"
	"github.com/submaline/services/gen/auth/v1/authv1connect"
	supervisorv1 "github.com/submaline/services/gen/supervisor/v1"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	talkv1 "github.com/submaline/services/gen/talk/v1"
	typesv1 "github.com/submaline/services/gen/types/v1"
	"github.com/submaline/services/logging"
	"github.com/submaline/services/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"strings"
)

type TalkServer struct {
	DB     *db.DBClient    // for mariadb
	Id     *snowflake.Node // for id generate
	Logger *zap.Logger     // for logging

	AuthClient *authv1connect.AuthServiceClient
	SvClient   *supervisorv1connect.SupervisorServiceClient
}

func (s *TalkServer) SendMessage(_ context.Context,
	req *connect.Request[talkv1.SendMessageRequest]) (
	*connect.Response[talkv1.SendMessageResponse], error) {
	senderUserId := req.Header().Get("X-Submaline-UserId")

	msg := req.Msg.Message
	msgId := fmt.Sprintf("ms|%s", xid.New().String())
	msg.Id = msgId
	msg.From = senderUserId // 強制付け替え

	// opを受け取るユーザーの生のuser_idが入る
	sendOpDest := []string{msg.From}
	recvOpDest := []string{msg.From}

	// 送信先のチェック
	switch {
	case strings.Contains(msg.To, "gr|"):
		// groupか?
		//opDestination
		// todo : impl
		return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("group is not implemented"))
	case strings.Contains(msg.To, "di|") && strings.Contains(msg.To, ".") && strings.Contains(msg.To, senderUserId):
		// 余計な部分を一回排除
		receiverUserId := strings.Replace(msg.To, "di|", "", 1)
		receiverUserId = strings.Replace(receiverUserId, ".", "", 1)
		receiverUserId = strings.Replace(receiverUserId, senderUserId, "", 1)

		// 自分宛ではないか
		if receiverUserId == senderUserId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid receiver"))
		}

		// 存在確認
		if !s.DB.IsAccountExists(receiverUserId) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid receiver"))
		}

		// opを相手が受け取れるように
		recvOpDest = append(recvOpDest, receiverUserId)

		// direct-chat id再構築
		msg.To = util.CreateDirectChatId(msg.From, receiverUserId)

	default:
		// 自分宛ではないか
		if msg.To == senderUserId {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid receiver"))
		}

		// 存在確認
		if !s.DB.IsAccountExists(msg.To) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid receiver"))
		}

		// opを相手が受け取れるように
		// msg.toはこの時点で相手のuser_idが入ってる
		recvOpDest = append(recvOpDest, msg.To)

		// direct-chat idに変更
		msg.To = util.CreateDirectChatId(msg.From, msg.To)
	}

	// 実際にdbに挿入したものを返してあげる
	resMsg, err := s.DB.CreateMessage(msg)
	if err != nil {
		// log
		if e_ := logging.ErrD(
			s.Logger,
			req.Spec().Procedure,
			err,
			"データベースにMessageを記録できませんでした",
			[]logging.DiscordRichMessageEmbedField{
				logging.GenerateDiscordRichMsgField("msgId", msgId, false),
			},
			os.Getenv("DISCORD_WEBHOOK_URL")); e_ != nil {
			log.Println(e_)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	// sv用のトークン生成
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
			Type:        typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE,
			Source:      msg.From,
			Destination: sendOpDest,
			Param1:      msg.Id,
			Param2:      msg.From,
			Param3:      msg.To,
			CreatedAt:   timestamppb.Now(),
		},
		{
			Id:          0,
			Type:        typesv1.OperationType_OPERATION_TYPE_SEND_MESSAGE_RECV,
			Source:      msg.From,
			Destination: recvOpDest,
			Param1:      msg.Id,
			Param2:      msg.From,
			Param3:      msg.To,
			CreatedAt:   timestamppb.Now(),
		},
	}})

	// トークンをくっつけてあげる
	recordReq.Header().Set("Authorization", fmt.Sprintf("Bearer %s", adminTokenResp.Msg.AuthToken.Token))
	// リクエスト送信
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
	// レスポンス作成
	res := connect.NewResponse(&talkv1.SendMessageResponse{Message: resMsg})

	return res, nil
}

func (s *TalkServer) SendReadReceipt(_ context.Context,
	_ *connect.Request[talkv1.SendReadReceiptRequest]) (
	*connect.Response[talkv1.SendReadReceiptResponse], error) {

	err := fmt.Errorf("unimplemented: SendReadReceipt")

	return nil, connect.NewError(connect.CodeUnimplemented, err)
}
