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

type Hash struct{}

func NewHash() *Hash {
	return &Hash{}
}

func (hash *Hash) GenerateNextLevelHashByPrevious(firstHah string, withSalt bool) string {
	var stringToHash string
	if withSalt {
		stringToHash = firstHah + salt
	} else {
		stringToHash = firstHah
	}
	h := sha256.New()
	h.Write([]byte(stringToHash))
	bs := h.Sum(nil)
	sha256hash := hex.EncodeToString(bs)

	return sha256hash
}

func (hash *Hash) GenerateFirstLevelHash() string {
	uuHash := uuid.New().String()
	h := sha256.New()
	h.Write([]byte(uuHash))
	bs := h.Sum(nil)
	sha256hashFirst := hex.EncodeToString(bs)

	return sha256hashFirst
}
