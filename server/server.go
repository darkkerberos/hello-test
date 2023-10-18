package server

import (
	"context"
	"log"

	"git.bluebird.id/logistic/commons/cert"
	"git.bluebird.id/logistic/commons/config"
	"git.bluebird.id/logistic/commons/constant"
	"git.bluebird.id/logistic/commons/logger"

	"git.bluebird.id/logistic/kit"
	"google.golang.org/grpc"

	"github.com/darkkerberos/hello-test/core"
	"github.com/darkkerberos/hello-test/handler"
	"github.com/darkkerberos/hello-test/postgres"
	"github.com/darkkerberos/hello-test/repository"
	"github.com/go-kit/kit/auth/jwt"

	apmgrpc "git.bluebird.id/lib/apm/grpc"
	hGrpc "github.com/darkkerberos/hello-test/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	// "github.com/go-kit/kit/auth/jwt"
)

type Server interface {
	Serve(addr string)
	Stop()
}

type grpcServer struct {
	kit      *kit.GRPCKit
	close    func()
	sayHello grpctransport.Handler
}

//NewGRPCServer ...
func NewGRPCServer() Server {
	cacert := config.Get(constant.CACertKey, "")
	serverCert := config.Get(constant.ServerCertKey, "")
	keyCert := config.Get(constant.KeyCertKey, "")
	var opts []grpc.ServerOption
	if cacert != "" && serverCert != "" && keyCert != "" {
		_, creds, err := cert.CreateTransportCredentials(cacert, serverCert, keyCert)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(apmgrpc.GetElasticAPMServerOptions(), grpc.Creds(*creds))

	} else {
		opts = apmgrpc.GetElasticAPMServerOptions()
	}
	//set database
	db, err := postgres.NewPostgresReaderWriter(repository.DBConfiguration{
		DBHost:     config.Get(constant.DBHostKey, "127.0.0.1"),
		DBName:     config.Get(constant.DBNameKey, ""),
		DBOptions:  config.Get(constant.DBOptionsKey, ""),
		DBPassword: config.Get(constant.DBPasswordKey, "password"),
		DBPort:     config.Get(constant.DBPortKey, "5432"),
		DBUser:     config.Get(constant.DBUserKey, "root"),
	})
	if err != nil {
		logger.Error("postgres connection err:", err.Error())
	}
	server := kit.NewGRPCServer(opts...)

	helloTestSvc := core.HelloTestService{DB: db}
	helloTestHandler := handler.NewHelloTestHandler(helloTestSvc)
	handlerOpts := []grpctransport.ServerOption{grpctransport.ServerBefore(jwt.GRPCToContext())}

	logger.Info("set helloTestService")

	return &grpcServer{
		kit: server,
		close: func() {
			db.Close()
		},
		sayHello: kit.NewGRPCHandler(helloTestHandler.SayHello, handlerOpts),
	}
}

func (g *grpcServer) Serve(addr string) {
	hGrpc.RegisterHelloTestServer(g.kit.Server, g)
	logger.Info("Running server on " + addr)
	g.kit.Run(addr)
}

func (g *grpcServer) Stop() {
	g.close()
	g.kit.GracefulStop()
}

func (g *grpcServer) SayHello(ctx context.Context, req *hGrpc.StringMessage) (*hGrpc.StringMessage, error) {
	_, resp, err := g.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return resp.(*hGrpc.StringMessage), nil
}
