package main

import (
	"QnA/internal/config"
	"QnA/internal/database/postgres"
	"QnA/internal/repository"
	"QnA/internal/server"
	"QnA/internal/transport/my_http/answer_handler"
	"QnA/internal/transport/my_http/question_handler"
	"QnA/internal/usecase"
	"log"
)

func main() {
	cfg := config.Load()

	gormDB, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// repos
	qRepo := repository.NewQuestionRepository(gormDB)
	aRepo := repository.NewAnswerRepository(gormDB)

	// usecases
	qUC := usecase.NewQuestionUsecase(qRepo)
	aUC := usecase.NewAnswerUsecase(qRepo, aRepo)

	// handlers
	qh := question_handler.NewQuestionHandler(qUC)
	ah := answer_handler.NewAnswerHandler(aUC)

	// router
	router := server.SetupRouter(qh, ah)

	// run
	server.Run(router, cfg.HTTPPort)
}
