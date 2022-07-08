package database

import (
	"database/sql"
	typesv1 "github.com/submaline/services/gen/protocol/types/v1"
)

// CreateAccount <Supervisor用> userIdとemailを使用してアカウント作成をします。
func (db *DBClient) CreateAccount(userId string, email string) (*typesv1.Account, error) {
	prep, err := db.Prepare("insert into accounts (user_id, email) values (?, ?)")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	_, err = prep.Exec(userId, email)
	if err != nil {
		return nil, err
	}

	return db.GetAccount(userId)
}

// GetAccount アカウント情報を取得します。
func (db *DBClient) GetAccount(userId string) (*typesv1.Account, error) {
	prep, err := db.Prepare("select email, user_name from accounts left join user_names on accounts.user_id = user_names.user_id where accounts.user_id = ?")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	var email string
	var userNameNullable sql.NullString
	err = prep.QueryRow(userId).Scan(&email, &userNameNullable)

	var userName string
	if userNameNullable.Valid {
		userName = userNameNullable.String
	}

	return &typesv1.Account{
		UserId:   userId,
		Email:    email,
		UserName: userName,
	}, err
}

func (db *DBClient) IsAccountExists(userId string) bool {
	_, err := db.GetAccount(userId)
	return err != nil
}
