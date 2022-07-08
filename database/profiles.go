package database

import (
	typesv1 "github.com/submaline/services/gen/protocol/types/v1"
	userv1 "github.com/submaline/services/gen/protocol/user/v1"
	"github.com/submaline/services/util"
	"strings"
)

func (db *DBClient) CreateProfile(userId string, displayName string, iconPath string) (*typesv1.Profile, error) {
	prep, err := db.Prepare("insert into profiles (user_id, display_name, icon_path) values (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	_, err = prep.Exec(userId, displayName, iconPath)
	if err != nil {
		return nil, err
	}

	return db.GetProfile(userId)
}

func (db *DBClient) GetProfile(userId string) (*typesv1.Profile, error) {
	prep, err := db.Prepare("select display_name, icon_path, status_message, metadata from profiles where user_id = ?")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	var displayName string
	var iconPath string
	var statusMessage string
	var metadata string
	err = prep.QueryRow(userId).Scan(&displayName, &iconPath, &statusMessage, &metadata)

	return &typesv1.Profile{
		UserId:        userId,
		DisplayName:   displayName,
		IconPath:      iconPath,
		StatusMessage: statusMessage,
		Metadata:      metadata,
	}, err
}

func (db *DBClient) UpdateProfile(userId string, request *userv1.UpdateProfileRequest) (*typesv1.Profile, error) {
	var values []interface{}
	var wheres []string

	if request.DisplayName != nil {
		wheres = append(wheres, "display_name = ?")
		values = append(values, request.GetDisplayName())
	}

	if request.IconPath != nil {
		wheres = append(wheres, "icon_path = ?")
		values = append(values, request.GetIconPath())
	}

	if request.StatusMessage != nil {
		wheres = append(wheres, "status_message = ?")
		values = append(values, request.GetStatusMessage())
	}

	if request.Metadata != nil {
		// 簡易jsonチェック.
		if !util.IsDBJson(request.GetMetadata()) {
			return nil, NewInvalidArgumentError("metadata")
		}
		wheres = append(wheres, "metadata = ?")
		values = append(values, request.GetMetadata())
	}

	// 式を用意
	prep, err := db.Prepare("update profiles set " + strings.Join(wheres, ", ") + " where user_id = ?")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	// (x, y...)はできるけど(x..., y)はできないぽいので、最後に加えてあげる
	values = append(values, userId)

	/// ...で配列の中身を展開
	_, err = prep.Exec(values...)
	if err != nil {
		return nil, err
	}

	// 新しくなった物を返してあげる.
	return db.GetProfile(userId)
}
