package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

// генерируем хэш первого порядка и делаем по нему хэш второго порядка + соль(его записываем в хранилище)
// - отправляем пользователю хэш первого порядка
// когда пользователь отправляет его назад(хэш первго порядка), то мы по нему генерируем хэш второго(с солью) и
// по нему находим в хранилище данные. В итоговую ссылку на скачивание будем добавлять хэш 1-го уровня и название архива

type Hash1lvl string
type Hash2lvl string

func GenerateNextLevelHashByPrevious(firstHah Hash1lvl, withSalt bool) Hash2lvl {
	var stringToHash string
	if withSalt {
		stringToHash = string(firstHah) + salt
	} else {
		stringToHash = string(firstHah)
	}
	sha256h := sha256.New()
	sha256h.Write([]byte(stringToHash))
	bs := sha256h.Sum(nil)
	sha256str := hex.EncodeToString(bs)

	return Hash2lvl(sha256str)
}

func GenerateFirstLevelHash() Hash1lvl {
	uuHash := uuid.New().String()
	sha256h := sha256.New()
	sha256h.Write([]byte(uuHash))
	bs := sha256h.Sum(nil)
	sha256str := hex.EncodeToString(bs)

	return Hash1lvl(sha256str)
}
