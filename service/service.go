package service

import "context"

type HelloTestService interface {
	SayHello(ctx context.Context, req string) (string, error)
}
