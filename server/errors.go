package server

import (
	"errors"
)

var ErrAdminOnly = errors.New("admin only")

const (
	ErrMsgFailedToGetUserDatFromFirebase = "Firebaseからユーザーデータを取得できませんでした"
)
