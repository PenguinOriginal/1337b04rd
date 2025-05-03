package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type SessionService interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error)
	DeleteExpiredSessions(ctx context.Context) error
	OverrideUserName(ctx context.Context, sessionID utils.UUID, newName string) error
}