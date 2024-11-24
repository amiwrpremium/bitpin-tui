package db

import (
	"gorm.io/gorm"
)

type Session struct {
	*gorm.Model
	ApiKey       string
	SecretKey    string
	AccessToken  string
	RefreshToken string
}

type Favorite struct {
	*gorm.Model
	Section string
	Symbol  string
	Count   int
}

type Setting struct {
	*gorm.Model
	Key   string
	Value string
}
