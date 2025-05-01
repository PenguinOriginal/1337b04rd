package handler


// Sort comments from comment_serv_impl.go, to distinguish between replies and regular 
// comments. Change html file accordingly

// type CommentThread struct {
// 	Comment *model.Comment
// 	Replies []*model.Comment
// }

// var threads []CommentThread
// repliesMap := make(map[utils.UUID][]*model.Comment)

// for _, c := range comments {
// 	if c.ParentCommentID != nil {
// 		repliesMap[*c.ParentCommentID] = append(repliesMap[*c.ParentCommentID], c)
// 	}
// }

// for _, c := range comments {
// 	if c.ParentCommentID == nil {
// 		thread := CommentThread{
// 			Comment: c,
// 			Replies: repliesMap[c.ID],
// 		}
// 		threads = append(threads, thread)
// 	}
// }

// Possible html file changes:
// {{ range .Threads }}
// 	<div class="top-level-comment">
// 		<p>{{ .Comment.Content }}</p>

// 		{{ range .Replies }}
// 			<div class="reply">
// 				<p>↳ {{ .Content }}</p>
// 			</div>
// 		{{ end }}
// 	</div>
// {{ end }}
