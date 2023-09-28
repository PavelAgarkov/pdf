package adapter

import (
	"errors"
	"pdf/internal"
	"pdf/internal/hash"
)

const (
	CookiesAlias = "archive"
)

type CookiesAdapter struct{}

func (ca *CookiesAdapter) GetAlias() string {
	return CookiesAlias
}

func NewCookiesAdapter() *CookiesAdapter {
	return &CookiesAdapter{}
}

func (ca *CookiesAdapter) IsAuthenticated(
	cookieInStorage internal.Hash2lvl,
	cookieValue internal.Hash1lvl,
	hit bool,
) (bool, error) {
	if cookieValue == "" {
		return false, errors.New("")
	}

	if !hit {
		return false, errors.New("")
	}

	hashFromCookie := hash.GenerateNextLevelHashByPrevious(cookieValue, true)

	if cookieInStorage != hashFromCookie {
		return false, errors.New("")
	}

	return true, nil
}
