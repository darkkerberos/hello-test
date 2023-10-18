package main

import (
	"fmt"

	"git.bluebird.id/logistic/commons/config"
	"git.bluebird.id/logistic/commons/constant"
	"github.com/darkkerberos/hello-test/server"
)

func main() {
	sv := server.NewGRPCServer()
	sv.Serve(fmt.Sprintf(":%s", config.Get(constant.ServicePortKey, "8085")))
	defer sv.Stop()
}
