package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	ColorWarn    = "#FFD60C" // yellow
	ColorErr     = "#FF0C0C" // red
	ColorSuccess = "#0CFF59" // green
	ColorInfo    = "#0C9EFF" // blue
)

type DiscordProfile struct {
	DisplayName string `json:"username"`
	Icon        string `json:"avatar_url"`
}

type DiscordRichMessage struct {
	DiscordProfile
	Content string                    `json:"content"`
	Embeds  []DiscordRichMessageEmbed `json:"embeds"`
}

type DiscordRichMessageEmbed struct {
	Title       string                         `json:"title"`
	Description string                         `json:"description"`
	Color       int64                          `json:"color"`
	Fields      []DiscordRichMessageEmbedField `json:"fields"`
	Author      DiscordRichMessageEmbedAuthor  `json:"author"`
}

type DiscordRichMessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type DiscordRichMessageEmbedAuthor struct {
	Name string `json:"name"`
}

func GenerateDiscordRichMsg(
	profile DiscordProfile,
	msg string,
	title string,
	desc string,
	color string,
	fields []DiscordRichMessageEmbedField,
	author string) DiscordRichMessage {
	color = strings.Replace(color, "#", "", 1)
	colorI, err := strconv.ParseInt(color, 16, 64)
	if err != nil {
		// #FFFFFF
		colorI = 16777215
	}
	d := DiscordRichMessage{
		DiscordProfile: profile,
		Content:        msg,
		Embeds: []DiscordRichMessageEmbed{
			{Title: title,
				Description: desc,
				Color:       colorI,
				Fields:      fields,
				Author:      DiscordRichMessageEmbedAuthor{Name: author}},
		},
	}
	return d
}

func GenerateDiscordRichMsgField(name string, value string, inline bool) DiscordRichMessageEmbedField {
	return DiscordRichMessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
}

func SendDiscordRichMessage(webhookUrl string, rich DiscordRichMessage) error {
	j, err := json.Marshal(rich)
	if err != nil {
		return err
	}
	log.Println(string(j))
	res, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(j))
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if 204 != res.StatusCode {
		return fmt.Errorf("failde to send message via discord: %v", string(resBody))
	}
	return nil
}
