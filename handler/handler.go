package handler

import (
	"git.bluebird.id/logistic/commons/logger"
	"git.bluebird.id/logistic/kit"
	pb "github.com/darkkerberos/hello-test/grpc"
	svc "github.com/darkkerberos/hello-test/service"
)

type HelloTestHandler struct {
	svc svc.HelloTestService
}

func NewHelloTestHandler(svc svc.HelloTestService) HelloTestHandler {
	return HelloTestHandler{
		svc: svc,
	}
}

func (h HelloTestHandler) SayHello(ctx kit.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.StringMessage)
	resp, err := h.svc.SayHello(ctx.Context, req.Value)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &pb.StringMessage{
		Value: resp,
	}, nil
}
