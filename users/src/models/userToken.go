package models

import "time"

type UserToken struct {
	Id        uint
	UserId    uint
	Token     string
	CreatedAt time.Time
	ExpiredAt time.Time
}
