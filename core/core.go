package core

import (
	"context"

	repo "github.com/darkkerberos/hello-test/repository"
	svc "github.com/darkkerberos/hello-test/service"
)

type helloTestService struct {
	db repo.PostgresReaderWriter
}

type HelloTestService struct {
	DB repo.PostgresReaderWriter
}

func NewHelloTestService(service HelloTestService) svc.HelloTestService {
	return &helloTestService{
		db: service.DB,
	}
}

func (h *helloTestService) SayHello(ctx context.Context, req string) (string, error) {
	msg := req
	return msg, nil
}
