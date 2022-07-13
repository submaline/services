package logging

import "go.uber.org/zap"

func Err(l *zap.Logger, source string, err error, msg string) {
	l.Error(msg,
		zap.String("source", source),
		zap.Error(err),
	)
}

func ErrD(l *zap.Logger, source string, err error, msg string, fields []DiscordRichMessageEmbedField, url string) error {
	Err(l, source, err, msg)

	prof := DiscordProfile{
		DisplayName: "Submaline/Log",
		Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
	}
	rich := GenerateDiscordRichMsg(
		prof,
		msg,
		"ERR",
		err.Error(),
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

func InfoD(l *zap.Logger, source string, msg string, fields []DiscordRichMessageEmbedField, url string) error {
	Info(l, source, msg)
	prof := DiscordProfile{
		DisplayName: "Submaline/Log",
		Icon:        "https://cdn.x0y14.workers.dev/250x250/e2890cb5-03d7-4176-ad0e-2071dec045fb",
	}

	rich := GenerateDiscordRichMsg(
		prof,
		msg,
		"INFO",
		"",
		ColorInfo,
		fields,
		source,
	)

	if err_ := SendDiscordRichMessage(url, rich); err_ != nil {
		return err_
	}

	return nil
}
