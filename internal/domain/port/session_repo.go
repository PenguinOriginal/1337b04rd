// Work on this later
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type SessionRepo interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error)
	DeleteExpiredSession(ctx context.Context) error
}
