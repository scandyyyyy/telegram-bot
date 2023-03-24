package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pocket "github.com/scandyyyyy/poket-sdk"
	"log"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/rep"
	"telegram-bot/pkg/rep/boltdb"
	"telegram-bot/pkg/server"
	"telegram-bot/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKye)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepositoriy(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	authServer := server.NewAuthServer(pocketClient, tokenRepository, cfg.TelegramBotURL)
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}

}
func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(rep.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(rep.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
