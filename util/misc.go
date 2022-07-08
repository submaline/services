package util

import (
	"encoding/json"
	"strings"
)

func IsJSON(str string) bool {
	// jsonの構文だけをみてるので、""が通ってしまう。
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsDBJson(str string) bool {
	hasCurly := strings.Contains(str, "{") && strings.Contains(str, "}")
	hasSquare := strings.Contains(str, "[") && strings.Contains(str, "]")
	return IsJSON(str) && (hasCurly || hasSquare)
}
