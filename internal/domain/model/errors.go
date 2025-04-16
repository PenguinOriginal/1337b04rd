package model

import "errors"

// We may not need all errors

// Common reusable errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input provided")
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden operation")
)

// Post-specific errors
var (
	ErrPostTitleRequired = errors.New("post title is required")
	ErrPostContentEmpty  = errors.New("post content cannot be empty")
	// ErrPostTooLong       = errors.New("post exceeds maximum allowed length")
	ErrPostNotFound = errors.New("post not found")
)

// Comment-specific errors
var (
	ErrCommentEmpty    = errors.New("comment cannot be empty")
	ErrCommentNotFound = errors.New("comment not found")
	ErrTooManyComments = errors.New("too many comments for this post")
)

// Session-related errors
var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionExpired   = errors.New("session expired")
	ErrInvalidSessionID = errors.New("invalid session ID")
)

// Misc
var (
	ErrDatabase       = errors.New("database error")
	ErrUUIDGeneration = errors.New("failed to generate UUID")
)
