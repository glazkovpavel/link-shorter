package link

import (
	"gorm.io/gorm"
	"math/rand"
)

type Link struct {
	gorm.Model
	URL  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		URL:  url,
		Hash: RandStringRunes(10),
	}
}

var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-+*&^%$#@!")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
