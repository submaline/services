package util

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
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

func ParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

func CreateDirectChatId(userId1 string, userId2 string) string {
	mem := []string{userId1, userId2}
	sort.Strings(mem)
	return fmt.Sprintf("di|%s", strings.Join(mem, "."))
}
