package db

import (
	"fmt"
	"strings"
)

func (db *DBClient) CreateOperationDestination(operationId int64, userIds []string) error {
	var questionPairs []string
	var valuePairs []any
	for _, userId := range userIds {
		questionPairs = append(questionPairs, "(?, ?)")
		valuePairs = append(valuePairs, operationId, userId)
	}

	query := fmt.Sprintf("insert into operation_destinations (operation_id, destination_user_id) values %s",
		strings.Join(questionPairs, ", "))

	prep, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = prep.Exec(valuePairs...)

	return err
}

func (db *DBClient) GetOperationDestinations(operationId int64) ([]string, error) {
	prep, err := db.Prepare("select destination_user_id from operation_destinations where operation_id = ?")
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query(operationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var destination []string
	for rows.Next() {
		var dest string
		if err := rows.Scan(&dest); err != nil {
			return destination, err
		}
		destination = append(destination, dest)
	}
	if err = rows.Err(); err != nil {
		return destination, err
	}

	return destination, nil
}
