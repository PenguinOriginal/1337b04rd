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
	ErrPostNotFound     = errors.New("post not found")
	ErrMissingTitle     = errors.New("post title is required")
	ErrMissingSessionID = errors.New("session ID is required")
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

// Triple-S related
var ErrBucketAlreadyExists = errors.New("bucket already exists")

// Misc
var (
	ErrDatabase       = errors.New("database error")
	ErrUUIDGeneration = errors.New("failed to generate UUID")
)
