package logging

import (
	"os"
	"testing"
)

func TestSendDiscordRichMessage(t *testing.T) {
	url := os.Getenv("DISCORD_WEBHOOK_URL")

	//field1 := GenerateDiscordRichMsgField("key", "value", false)
	msg := GenerateDiscordRichMsg(
		DiscordProfile{
			DisplayName: "Submaline",
			Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
		},
		"test message",
		"the title",
		"this is test",
		ColorWarn,
		nil,
		"TestSendDiscordRichMessage",
	)

	if err := SendDiscordRichMessage(url, msg); err != nil {
		t.Fatalf("failed to send msg: %v", err)
	}

}
