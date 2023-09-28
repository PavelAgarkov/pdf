package service

import (
	"errors"
	"pdf/internal"
	"pdf/internal/hash"
	"strings"
)

func IsAuthenticated(
	storageHash internal.Hash2lvl,
	requestHash internal.Hash1lvl,
) (bool, error) {
	if requestHash == "" {
		return false, errors.New("")
	}

	hashFromRequest := hash.GenerateNextLevelHashByPrevious(requestHash, true)

	if storageHash != hashFromRequest {
		return false, errors.New("")
	}

	return true, nil
}

func ParseBearerHeader(bearer string) string {
	if strings.Contains(bearer, internal.Bearer) {
		token := strings.Split(bearer, internal.BearerSeparator)
		if len(token) == 2 && token[0] == internal.Bearer {
			return token[1]
		}
		return ""
	}
	return ""
}

func GenerateBearerToken() string {
	return string(hash.GenerateFirstLevelHash())
}
