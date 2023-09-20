package internal

import (
	storage2 "pdf/internal/storage"
	"time"
)

type UserData struct {
	hash1lvl  storage2.Hash1lvl // это и будет ключ для записи в куки
	hash2lvl  storage2.Hash2lvl // это и будет ключ для записи в основное хранилище и для ссылки
	expiredAt time.Time
}

func NewUserData(
	hash1lvl storage2.Hash1lvl,
	hash2lvl storage2.Hash2lvl,
	expiredAt time.Time,
) *UserData {
	return &UserData{
		hash1lvl:  hash1lvl,
		hash2lvl:  hash2lvl,
		expiredAt: expiredAt,
	}
}

func (ud *UserData) GetExpiredAt() time.Time {
	return ud.expiredAt
}

func (ud *UserData) GetHash1Lvl() storage2.Hash1lvl {
	return ud.hash1lvl
}

func (ud *UserData) GetHash2Lvl() storage2.Hash2lvl {
	return ud.hash2lvl
}
