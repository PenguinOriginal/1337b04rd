package model

import "time"

type Post struct {
	PostID     UUID
	SessionID  UUID
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
