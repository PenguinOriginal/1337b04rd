package model

import (
	"1337b04rd/pkg/utils"
	"strings"
	"time"
)

type Post struct {
	PostID     utils.UUID
	SessionID  utils.UUID
	UserName   string
	Title      string
	Content    string
	ImageURLs  []string
	CreatedAt  time.Time
	IsArchived bool
}

// Work on this later
func (p *Post) IsExpired(lastCommentAt *time.Time) bool {
	if lastCommentAt == nil {
		return time.Since(p.CreatedAt) > 10*time.Minute
	}
	return time.Since(*lastCommentAt) > 15*time.Minute
}

func (p *Post) ValidatePost() error {

	if strings.TrimSpace(p.Title) == "" {
		return ErrMissingTitle
	}
	if p.SessionID == "" {
		return ErrMissingSessionID
	}
	// if there's no username, db sets it as "Anonymous"
	return nil
}
