package logging

import "go.uber.org/zap"

func Err(l *zap.Logger, source string, err error, msg string) {
	l.Error(msg,
		zap.String("source", source),
		zap.Error(err),
	)
}

func ErrD(l *zap.Logger, source string, err error, desc string, fields []DiscordRichMessageEmbedField, url string) error {
	Err(l, source, err, desc)

	prof := DiscordProfile{
		DisplayName: "Submaline/Log",
		Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
	}

	fields = append(fields, GenerateDiscordRichMsgField("detail", err.Error(), false))

	rich := GenerateDiscordRichMsg(
		prof,
		"",
		"ERR",
		desc,
		ColorErr,
		fields,
		source,
	)

	if err_ := SendDiscordRichMessage(url, rich); err_ != nil {
		return err_
	}

	return nil
}

func Info(l *zap.Logger, source, msg string) {
	l.Info(msg,
		zap.String("source", source),
	)
}

func InfoD(l *zap.Logger, source string, desc string, fields []DiscordRichMessageEmbedField, url string) error {
	Info(l, source, desc)
	//prof := DiscordProfile{
	//	DisplayName: "Submaline/Log",
	//	Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
	//}
	//
	//rich := GenerateDiscordRichMsg(
	//	prof,
	//	"",
	//	"INFO",
	//	desc,
	//	ColorInfo,
	//	fields,
	//	source,
	//)
	//
	//if err_ := SendDiscordRichMessage(url, rich); err_ != nil {
	//	return err_
	//}

	return nil
}
