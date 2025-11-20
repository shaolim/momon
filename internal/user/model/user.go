package model

import (
	"errors"
	"time"
)

type User struct {
	ID          int64
	LineUserID  string
	DisplayName string
	Status      UserStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) Validate() error {
	if u.LineUserID == "" {
		return errors.New("line user id must not be empty")
	}

	return nil
}

type UserStatus string

const (
	UserStatusActive   = "ACTIVE"
	UserStatusInActive = "INACTIVE"
)
