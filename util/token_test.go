package util

import (
	"log"
	"testing"
)

func TestGenToken(t *testing.T) {
	res, err := GenToken("", "")
	if err != nil {
		t.Fatalf("failed to gen token: %v", err)
	}

	log.Println(res)
}

func TestGenTokenWithRefresh(t *testing.T) {

}
