package model

type User struct {
	ID          string
	LineUserID  string
	DisplayName string
	Status      UserStatus
}

type UserStatus string

const (
	UserStatusActive   = "ACTIVE"
	UserStatusInActive = "INACTIVE"
)
