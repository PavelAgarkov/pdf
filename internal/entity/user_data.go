package entity

import (
	"pdf/internal"
	"time"
)

type UserData struct {
	hash1lvl  internal.Hash1lvl // это и будет ключ для записи в куки
	hash2lvl  internal.Hash2lvl // это и будет ключ для записи в основное хранилище и для ссылки
	expiredAt time.Time
}

func NewUserData(
	hash1lvl internal.Hash1lvl,
	hash2lvl internal.Hash2lvl,
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

func (ud *UserData) GetHash1Lvl() internal.Hash1lvl {
	return ud.hash1lvl
}

func (ud *UserData) GetHash2Lvl() internal.Hash2lvl {
	return ud.hash2lvl
}
