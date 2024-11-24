package db

func UpsertSession(apiKey, secretKey, accessToken, refreshToken string) {
	var session Session
	DB.Where(&Session{ApiKey: apiKey}).Or(&Session{SecretKey: secretKey}).Or(&Session{AccessToken: accessToken}).Or(&Session{RefreshToken: refreshToken}).First(&session)

	session.ApiKey = apiKey
	session.SecretKey = secretKey
	session.AccessToken = accessToken
	session.RefreshToken = refreshToken
	DB.Save(&session)
}

func GetSession() *Session {
	var session Session
	err := DB.First(&session).Error

	if err != nil {
		return &Session{}
	}

	return &session
}

func DeleteSession() error {
	return DB.Unscoped().Where("1 = 1").Delete(&Session{}).Error
}

func UpsertFavorite(section, symbol string) {
	var favorite Favorite
	result := DB.Where(&Favorite{Section: section, Symbol: symbol}).First(&favorite)

	if result != nil && result.Error != nil {
		DB.Create(&Favorite{
			Symbol:  symbol,
			Section: section,
			Count:   1,
		})
	} else {
		favorite.Count++
		DB.Save(&favorite)
	}
}

func GetFavorites(section string) []string {
	var favorites []Favorite
	result := DB.
		Select("symbol").
		Where(&Favorite{Section: section}).
		Order("count DESC").
		Limit(10).
		Find(&favorites)

	if result != nil && result.Error != nil {
		return []string{}
	}

	var symbols []string
	for _, favorite := range favorites {
		symbols = append(symbols, favorite.Symbol)
	}

	return symbols
}
