package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Auth struct {
	secretKey string
}

func NewAuth(secretKey string) *Auth {
	return &Auth{secretKey: secretKey}
}

func (a *Auth) IsValidToken(token string) bool {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false
	}

	headerPayload := parts[0] + "." + parts[1]
	signature := parts[2]

	mac := hmac.New(sha256.New, []byte(a.secretKey))
	mac.Write([]byte(headerPayload))
	expectedSignature := mac.Sum(nil)
	expectedSignatureBase64 := base64.RawURLEncoding.EncodeToString(expectedSignature)

	if expectedSignatureBase64 != signature {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var payload struct {
		Nbf int64 `json:"nbf"`
		Exp int64 `json:"exp"`
	}

	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return false
	}

	now := time.Now().Unix()

	if now < payload.Nbf || now > payload.Exp {
		return false
	}

	return true
}
