// Work on this later
package port

import "1337b04rd/internal/domain/model"

type CommentService interface {
	CreateComment(comment *model.Comment) error
	// others later

}
