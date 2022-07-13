package logging

import (
	"os"
	"testing"
)

const (
	ColorWarn    = "#FFD60C" // yellow
	ColorErr     = "#FF0C0C" // red
	ColorSuccess = "#0CFF59" // green
	ColorInfo    = "#0C9EFF" // blue
)

func TestSendDiscordRichMessage(t *testing.T) {
	url := os.Getenv("DISCORD_WEBHOOK_URL")

	field1 := GenerateDiscordRichMsgField("key", "value", false)
	msg := GenerateDiscordRichMsg(
		DiscordProfile{
			DisplayName: "Submaline",
			Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
		},
		"test message",
		"the title",
		"this is test",
		ColorWarn,
		[]DiscordRichMessageEmbedField{field1},
		"TestSendDiscordRichMessage",
	)

	if err := SendDiscordRichMessage(url, msg); err != nil {
		t.Fatal(err)
	}

}
