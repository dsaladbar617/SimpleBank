package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPayload(username string, duration time.Duration) (*Payload, error){
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: tokenId,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

