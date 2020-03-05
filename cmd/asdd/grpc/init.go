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

// ██████╗ ██████╗ ██████╗  ██████╗    ███╗   ███╗███████╗████████╗██╗  ██╗ ██████╗ ██████╗ ███████╗
//██╔════╝ ██╔══██╗██╔══██╗██╔════╝    ████╗ ████║██╔════╝╚══██╔══╝██║  ██║██╔═══██╗██╔══██╗██╔════╝
//██║  ███╗██████╔╝██████╔╝██║         ██╔████╔██║█████╗     ██║   ███████║██║   ██║██║  ██║███████╗
//██║   ██║██╔══██╗██╔═══╝ ██║         ██║╚██╔╝██║██╔══╝     ██║   ██╔══██║██║   ██║██║  ██║╚════██║
//╚██████╔╝██║  ██║██║     ╚██████╗    ██║ ╚═╝ ██║███████╗   ██║   ██║  ██║╚██████╔╝██████╔╝███████║
//╚═════╝ ╚═╝  ╚═╝╚═╝      ╚═════╝    ╚═╝     ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝

func (s *Server) Init(ctx context.Context, in *api.Pool) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Init")
	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Create(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")
	var apiOutcome api.Outcome

	// create a new datased, children of main dataset asd
	// and populate with a template unit

	return &apiOutcome, nil
}

func (s *Server) Destroy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Destroy")

	var apiOutcome api.Outcome

	// destroy a solution datased, with ALL backups, snapshots

	return &apiOutcome, nil
}

func (s *Server) Backup(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Backup")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Restore(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Restore")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Deploy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Deploy")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Retreat(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Retreat")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Start(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Start")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Stop(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Stop")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Expose(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Expose")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Contain(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Contain")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Snapshot(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Snapshot")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Rollback(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Rollback")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
