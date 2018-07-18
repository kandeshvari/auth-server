package server

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/kandeshvari/gin-jwt-middleware"
	"github.com/satori/go.uuid"
)

type RefreshTokenStorage struct {
	jwt.IRefreshTokenStorage
	rw      *sync.RWMutex
	revoked map[string]int64
	db      *sqlx.DB
}

type Session struct {
	Expire int64 `sql:"expire"`
	Token  int64 `sql:"token"`
}

func NewRefreshTokenStorage(driver, dbconnect string) (*RefreshTokenStorage, error) {
	return &RefreshTokenStorage{
		revoked: make(map[string]int64),
		rw:      &sync.RWMutex{},
		db:      sqlx.MustConnect(driver, dbconnect),
	}, nil
}

// Generate refresh token
func (ts *RefreshTokenStorage) Issue() (string, error) {
	return uuid.NewV4().String(), nil
}

// Check is refresh token already expired
func (ts *RefreshTokenStorage) IsExpired(token string) bool {
	var result Session
	query := "SELECT expire FROM sessions WHERE token = ?"
	err := ts.db.QueryRowx(query, token).StructScan(&result)
	if err != nil {
		log.Printf("ERROR: IsExpire(): %s", err)
		return true
	}
	log.Printf("IsExpire(): SELECT OK: %#v", result)
	if result.Expire < time.Now().Unix() {
		log.Printf("IsExpire(): token expired")
		return true
	}

	log.Printf("IsExpire(): token not expired")
	return false
}

// Add/Update refresh token in storage
func (ts *RefreshTokenStorage) Update(token string, refreshTimeout time.Duration, payload map[string]interface{}, c *gin.Context) error {
	expire := time.Now().Add(refreshTimeout).Unix()
	queryInsert := "INSERT IGNORE INTO sessions (token, expire, user_id, ip, agent, os) VALUES (?, ?, ?, ?, ?, ?)"
	queryUpdate := "UPDATE sessions SET expire = ?, ip = ?, agent = ?, os = ? WHERE token = ?"
	userId := payload["user_id"]
	ip := c.ClientIP()
	agent := "Unknown"
	os := "Unknown"
	_, err := ts.db.Exec(queryInsert, token, expire, userId, ip, agent, os)
	if err != nil {
		log.Printf("insert failed: %s", err)
		return err
	}
	_, err = ts.db.Exec(queryUpdate, expire, ip, agent, os, token)
	if err != nil {
		log.Printf("update failed: %s", err)
		return err
	}
	return nil
}

// Delete refresh token from storage
func (ts *RefreshTokenStorage) Delete(token string) error {
	query := "DELETE FROM sessions WHERE token = ?"
	_, err := ts.db.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

// Revoke refresh token
func (ts *RefreshTokenStorage) Revoke(token string, accessTokenTimeout time.Duration) error {
	ts.rw.Lock()
	ts.revoked[token] = time.Now().Add(accessTokenTimeout).Unix()
	ts.rw.Unlock()

	return nil
}

// Check is refresh token was revoked
func (ts *RefreshTokenStorage) IsRevoked(token string) bool {
	ts.rw.RLock()
	_, ok := ts.revoked[token]
	ts.rw.RUnlock()
	return ok
}
