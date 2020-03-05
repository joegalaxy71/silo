package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
}

func Init(errGrpc chan<- error) {

	_log := helpers.InitLogs(true)
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		err = errors.Wrap(err, "trying to listen on tcp")
		errGrpc <- err
		return
	}
	// create a server instance
	s := Server{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	api.RegisterAsdServer(grpcServer, &s)
	// start the server
	_log.Debugf("listening for grpc connections on port: 7777")
	if err := grpcServer.Serve(lis); err != nil {
		err = errors.Wrap(err, "failed to serve gRPC")
		errGrpc <- err
		return
	}
}
