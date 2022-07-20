package util

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 生成されたトークンを保存する
var tokenCache map[string]TokenData

// 初期化してあげる
func init() {
	tokenCache = map[string]TokenData{}
}

type TokenData struct {
	IdToken   string
	Refresh   string
	ExpiresAt time.Time
	ExpiresIn string // 互換
	UID       string // 互換
}

// SetAdminClaim shを書くのがめんどくさかった。アドミン用のカスタムクレームをつけるスクリプト
func SetAdminClaim(uid string) error {
	// Firebase appの初期化
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// Set admin privilege on the user corresponding to uid.
	claims := map[string]interface{}{"admin": true}
	err = client.SetCustomUserClaims(context.Background(), uid, claims)
	return err
}

type genTokenRequest struct {
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	ReturnSecureToken bool   `json:"return_secure_token,omitempty"`
}

type genTokenResponse struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"IdToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"ExpiresIn"`
}

// todo : error handle
type genTokenError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
		Status  string `json:"status"`
		Details []struct {
			Type     string `json:"@type"`
			Reason   string `json:"reason"`
			Domain   string `json:"domain"`
			Metadata struct {
				Service string `json:"service"`
			} `json:"metadata"`
		} `json:"details"`
	} `json:"error"`
}

type genTokenWithRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type genTokenWithRefreshResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	UserId       string `json:"user_id"`
	ProjectId    string `json:"project_id"`
}

// todo : error handle
type genTokenWithRefreshError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func GenToken(email, password string) (*TokenData, error) {
	url_ := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s",
		os.Getenv("FIREBASE_WEB_API_KEY"))
	bin := genTokenRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	dataBin, err := json.Marshal(bin)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url_, "application/json", bytes.NewBuffer(dataBin))
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result genTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	expiresIn, err := strconv.ParseInt(result.ExpiresIn, 10, 64)
	if err != nil {
		var resultErr genTokenError
		_ = json.Unmarshal(body, &resultErr)
		return nil, errors.Wrap(err, resultErr.Error.Message)
	}

	now := time.Now()
	expiresAt := now.Add(time.Second * time.Duration(expiresIn))
	return &TokenData{
		IdToken:   result.IdToken,
		Refresh:   result.RefreshToken,
		ExpiresAt: expiresAt,
		ExpiresIn: result.ExpiresIn,
		UID:       result.LocalId,
	}, nil
}

func GenTokenWithRefresh(refreshToken string) (*TokenData, error) {
	url_ := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s",
		os.Getenv("FIREBASE_WEB_API_KEY"))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	res, err := http.Post(url_, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	// todo : remove
	log.Println(string(body))

	if err != nil {
		return nil, err
	}

	var result genTokenWithRefreshResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	expiresIn, err := strconv.ParseInt(result.ExpiresIn, 10, 64)
	if err != nil {

		return nil, err
	}
	now := time.Now()
	expiresAt := now.Add(time.Second * time.Duration(expiresIn))

	return &TokenData{
		IdToken:   result.IdToken,
		Refresh:   result.RefreshToken,
		ExpiresAt: expiresAt,
		ExpiresIn: result.ExpiresIn,
		UID:       result.UserId,
	}, nil
}

func GenerateToken(email string, password string, renew bool) (*TokenData, error) {
	if os.Getenv("FIREBASE_WEB_API_KEY") == "" {
		if err := godotenv.Load(".env"); err != nil {
			return nil, errors.Wrap(err, "so FIREBASE_WEB_API_KEY was empty, i tried load .env but failed")
		}
	}

	// Q.再生成を強制しますか?
	if !renew {
		// A.いいえ、強制しません
		cache, ok := tokenCache[email]
		// キャッシュが存在する
		if ok {
			// 現在の時刻がトークンの期限切れの前か?
			if time.Now().Before(cache.ExpiresAt) {
				// これは使えます
				log.Printf("<GenerateToken> USE CACHE: %v\n", email)
				return &cache, nil
			} else {
				// リフレッシュトークンを使用して再発行
				tData, err := GenTokenWithRefresh(cache.Refresh)
				if err != nil {
					return nil, err
				}
				// 再生成成功
				// キャッシュとして記憶してあげる
				tokenCache[email] = *tData
				log.Printf("<GenerateToken> REFRESH & CACHE: %v\n", email)
				return tData, nil
			}
		}
	}
	// A.はい、強制します (& キャッシュが使えなかったので生成します)
	tData, err := GenToken(email, password)
	if err != nil {
		return nil, err
	}
	// キャッシュとして記憶
	tokenCache[email] = *tData

	log.Printf("<GenerateToken> GENERATE & CACHE: %v\n", email)
	return tData, nil
}
