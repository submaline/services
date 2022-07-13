package logging

import (
	"strconv"
	"strings"
)

type DiscordRichMessage struct {
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
		Content: msg,
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
