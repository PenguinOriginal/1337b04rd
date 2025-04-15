package model

import "time"

type Session struct {
	SessionID UUID
	AvatarURL string
	CreatedAt time.Time
	ExpiresAt time.Time
}
