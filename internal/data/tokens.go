package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserId    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  ScopeActivation,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base64.StdEncoding.EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (t TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.insert(token)
	return token, err
}

func (t TokenModel) insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
			  VALUES ($1, $2, $3, $4)`

	args := []interface{}{
		token.Hash,
		token.UserId,
		token.Expiry,
		token.Scope,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, args...)
	return err
}

func (t TokenModel) DeleteAllForUser(userId int64, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{userId, scope}
	_, err := t.DB.ExecContext(ctx, query, args...)
	return err
}

type TokenRepository interface {
	New(userId int64, ttl time.Duration, scope string) (*Token, error)
}
