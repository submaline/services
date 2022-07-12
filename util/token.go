package util

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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

type GenTokenRequest struct {
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	ReturnSecureToken bool   `json:"return_secure_token,omitempty"`
}

type GenTokenResponse struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

func GenerateToken(email string, password string) (*GenTokenResponse, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s",
		os.Getenv("FIREBASE_WEB_API_KEY"))
	bin := GenTokenRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	dataBin, err := json.Marshal(bin)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(dataBin))
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result GenTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}
