package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber int       `json:"phonenumber"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific phoneNumber and duration
func NewPayload(phoneNumber int, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenId,
		PhoneNumber: phoneNumber,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}