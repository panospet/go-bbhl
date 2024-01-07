package main

import (
	"log"

	"github.com/caarlos0/env/v9"

	"go-bbhl/uploaders/telegram"
)

type config struct {
	TelegramToken string `env:"TELE_TOKEN,required"`
	ChatId        int64  `env:"CHAT_ID,required"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	upl, err := telegram.NewTelegramUploader(
		cfg.TelegramToken,
		cfg.ChatId,
		// todo change to env var
		"http://65.109.174.137:8081/bot%s/%s",
	)
	if err != nil {
		log.Fatalf("Error creating uploader: %v", err)
	}

	log.Printf("uploading video to telegram channel...")
	if err := upl.UploadVideo("~/gopher.mp4", "sample test"); err != nil {
		log.Fatalf("Error uploading video: %v", err)
	}
}
