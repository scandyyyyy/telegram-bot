package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"strconv"
	"telegram-bot/pkg/rep"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepositoriy(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}
func (r *TokenRepository) Save(chatID int64, token string, bucket rep.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))

	})
}
func (r *TokenRepository) Get(chatID int64, bucket rep.Bucket) (string, error) {
	var token string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func intToBytes(a int64) []byte {
	return []byte(strconv.FormatInt(a, 10))
}