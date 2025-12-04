package server

import (
	"QnA/internal/transport/my_http/answer_handler"
	"QnA/internal/transport/my_http/question_handler"
	"net/http"
)

func SetupRouter(qh *question_handler.QuestionHandler, ah *answer_handler.AnswerHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /questions", qh.Create)
	mux.HandleFunc("GET /questions", qh.GetAll)
	mux.HandleFunc("GET /questions/{id}", qh.GetByID)
	mux.HandleFunc("DELETE /questions/{id}", qh.Delete)

	mux.HandleFunc("POST /questions/{id}/answers", ah.Create)
	mux.HandleFunc("GET /answers/{id}", ah.GetByID)
	mux.HandleFunc("DELETE /answers/{id}", ah.Delete)

	return mux
}
