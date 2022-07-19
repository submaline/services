package util

import (
	"log"
	"testing"
)

const (
	email    = ""
	password = ""
)

func TestGenToken(t *testing.T) {
	res, err := GenToken(email, password)
	if err != nil {
		t.Fatalf("failed to gen token: %v", err)
	}

	log.Println(res)
}

func TestGenTokenWithRefresh(t *testing.T) {
	res, err := GenToken(email, password)
	if err != nil {
		t.Fatalf("failed to gen token: %v", err)
	}

	tok, err := GenTokenWithRefresh(res.Refresh)
	if err != nil {
		t.Fatalf("failed to refresh: %v", err)
	}

	log.Println(tok)

}
