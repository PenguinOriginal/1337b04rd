package port

import "1337b04rd/internal/domain/model"

type SessionService interface {
	CreateSession(session *model.Session) error
	// others later
}
