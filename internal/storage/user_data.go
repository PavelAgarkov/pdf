package storage

import (
	"time"
)

type UserData struct {
	hash1lvl  Hash1lvl // это и будет ключ для записи в куки
	hash2lvl  Hash2lvl // это и будет ключ для записи в основное хранилище и для ссылки
	expiredAt time.Time
}

func NewUserData(
	hash1lvl Hash1lvl,
	hash2lvl Hash2lvl,
	expiredAt time.Time,
) *UserData {
	return &UserData{
		hash1lvl:  hash1lvl,
		hash2lvl:  hash2lvl,
		expiredAt: expiredAt,
	}
}
