package telegram

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramUploader struct {
	Token          string
	ChatId         int64
	BotApiEndpoint string
	bot            *tgbotapi.BotAPI
}

func NewTelegramUploader(
	token string,
	chatId int64,
	botApiEndpoint string,
) (*TelegramUploader, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(
		token,
		botApiEndpoint,
		&http.Client{},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %v", err)
	}
	return &TelegramUploader{
		Token:          token,
		ChatId:         chatId,
		BotApiEndpoint: botApiEndpoint,
		bot:            bot,
	}, nil
}

func (o *TelegramUploader) UploadVideo(
	path string,
	caption string,
) error {
	vid := tgbotapi.NewVideo(
		o.ChatId,
		tgbotapi.FilePath(path),
	)
	vid.Caption = caption

	_, err := o.bot.Send(vid)
	if err != nil {
		return fmt.Errorf("error sending video: %v", err)
	}

	return nil
}
