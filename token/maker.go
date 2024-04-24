package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific phoneNumber and duration
	CreateToken(phoneNumber int, duration time.Duration) (string, *Payload, error)
	// VerifyToken check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
