package main

import (
	"context"
	"github.com/kaverhovsky/chat/internal/domain/repository"
	"github.com/kaverhovsky/chat/internal/domain/usecase"
	"log"
	"time"
)

func main() {
	// TODO logger
	pgrepo := repository.NewPostgresRepository()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO config
	if err := pgrepo.Init(ctx, "postgres://postgres:password@localhost:5432/pechat_users"); err != nil {
		// TODO log error and exit gracefully
		log.Fatal("err")
	}

	_ = usecase.NewUseCase(pgrepo)
}
