package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	TokenType       string   `json:"token_type"`
	Exp             int      `json:"exp"`
	Jti             string   `json:"jti"`
	UserId          int      `json:"user_id"`
	Ip              []string `json:"ip"`
	ApiCredentialId int      `json:"api_credential_id"`
}

func (j JWT) HumanReadable() string {
	return fmt.Sprintf("Token Type: %s\nExp: %d\nJti: %s\nUserId: %d\nIp: %v\nApi Credential Id: %d", j.TokenType, j.Exp, j.Jti, j.UserId, j.Ip, j.ApiCredentialId)
}

func (j JWT) IsExpired() bool {
	return j.Exp < int(time.Now().Unix())
}

func (j JWT) IsExpiredIn(t time.Duration) bool {
	return j.Exp < int(time.Now().Add(t).Unix())
}

func DecodeJWT(tokenString string) (*JWT, error) {
	// Parse the JWT token.
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract the claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims as MapClaims")
	}

	// Marshal claims back to JSON.
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	// Unmarshal JSON into the T struct.
	var payload JWT
	if err := json.Unmarshal(claimsJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims into struct: %w", err)
	}

	return &payload, nil
}
