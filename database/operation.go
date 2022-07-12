package database

import (
	typesv1 "github.com/submaline/services/gen/types/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (db *DBClient) CreateOperation(
	operationId int64,
	operationType typesv1.OperationType,
	source string,
	p1 string,
	p2 string,
	p3 string,
	createdAt time.Time) (*typesv1.Operation, error) {
	prep, err := db.Prepare("insert into operations (id, type, source, param1, param2, param3, created_at) value (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	_, err = prep.Exec(operationId, operationType.Number(), source, p1, p2, p3, createdAt)

	return nil, err
}

func (db *DBClient) GetOperationWithOperationId(operationId int64) (*typesv1.Operation, error) {
	prep, err := db.Prepare("select type, source, param1, param2, param3, created_at from operations where id = ?")
	if err != nil {
		return nil, err
	}

	var opType int32
	var source string
	var param1 string
	var param2 string
	var param3 string
	var createdAt time.Time
	err = prep.QueryRow(operationId).Scan(&opType, &source, &param1, &param2, &param3, &createdAt)
	if err != nil {
		return nil, err
	}

	destinations, err := db.GetOperationDestinations(operationId)
	if err != nil {
		return nil, err
	}

	return &typesv1.Operation{
		Id:          operationId,
		Type:        typesv1.OperationType(opType),
		Source:      source,
		Destination: destinations,
		Param1:      param1,
		Param2:      param2,
		Param3:      param3,
		CratedAt:    timestamppb.New(createdAt),
	}, nil
}
