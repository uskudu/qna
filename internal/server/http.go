package server

import (
	"log"
	"net/http"
)

func Run(router http.Handler, port string) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Println("listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
