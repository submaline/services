package database

import (
	typesv1 "github.com/submaline/services/gen/protocol/types/v1"
	"github.com/submaline/services/util"
)

func (db *DBClient) CreateMessage(msg *typesv1.Message) (*typesv1.Message, error) {
	if !util.IsDBJson(msg.Metadata) {
		return nil, NewInvalidArgumentError("metadata")
	}

	prep, err := db.Prepare("insert into messages (id, `from`, `to`, content_type, text, metadata) value (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	_, err = prep.Exec(msg.Id, msg.From, msg.To, msg.ContentType, msg.Text, msg.Metadata)
	if err != nil {
		return nil, err
	}

	return db.GetMessageWithMessageId(msg.Id)
}

func (db *DBClient) GetMessageWithMessageId(messageId string) (*typesv1.Message, error) {
	prep, err := db.Prepare("select `from`, `to`, content_type, text, metadata from messages where id = ?")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(messageId)

	var from string
	var to string
	var contentType int32
	var text string
	var metadata string

	err = row.Scan(&from, &to, &contentType, &text, &metadata)
	if err != nil {
		return nil, err
	}

	return &typesv1.Message{
		Id:          messageId,
		From:        from,
		To:          to,
		ContentType: typesv1.MessageContentType(contentType),
		Text:        text,
		Metadata:    metadata,
	}, nil
}
