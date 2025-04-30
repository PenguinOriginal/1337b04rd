package model

import (
	"1337b04rd/pkg/utils"
	"time"
)

type Session struct {
	SessionID utils.UUID
	AvatarURL string
	CreatedAt time.Time
	ExpiresAt time.Time
}
