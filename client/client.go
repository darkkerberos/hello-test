package client

import (
	"context"

	"git.bluebird.id/logistic/kit"
	pb "github.com/darkkerberos/hello-test/grpc"
	"github.com/darkkerberos/hello-test/service"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

const (
	grpcHello = "grpc.HelloTest"
)

type clientService struct {
	sayHello endpoint.Endpoint
}

//NewHelloTestService ...
func NewHelloTestService(cc *grpc.ClientConn, opts ...grpctransport.ClientOption) service.HelloTestService {
	return clientService{
		sayHello: kit.NewGRPCClientHandler(
			cc,
			grpcHello,
			"SayHello",
			pb.StringMessage{},
			opts...),
	}
}

func (c clientService) SayHello(ctx context.Context, request string) (string, error) {
	resp, err := c.sayHello(ctx, &pb.StringMessage{
		Value: request,
	})
	if err != nil {
		return "", err
	}
	return resp.(*pb.StringMessage).Value, nil
}
