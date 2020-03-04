package main

import (
	"asd/common/api"
	"context"
)

// ██████╗ ██████╗ ██████╗  ██████╗    ███╗   ███╗███████╗████████╗██╗  ██╗ ██████╗ ██████╗ ███████╗
//██╔════╝ ██╔══██╗██╔══██╗██╔════╝    ████╗ ████║██╔════╝╚══██╔══╝██║  ██║██╔═══██╗██╔══██╗██╔════╝
//██║  ███╗██████╔╝██████╔╝██║         ██╔████╔██║█████╗     ██║   ███████║██║   ██║██║  ██║███████╗
//██║   ██║██╔══██╗██╔═══╝ ██║         ██║╚██╔╝██║██╔══╝     ██║   ██╔══██║██║   ██║██║  ██║╚════██║
//╚██████╔╝██║  ██║██║     ╚██████╗    ██║ ╚═╝ ██║███████╗   ██║   ██║  ██║╚██████╔╝██████╔╝███████║
//╚═════╝ ╚═╝  ╚═╝╚═╝      ╚═════╝    ╚═╝     ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝

func (s *Server) Init(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Create(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	var apiOutcome api.Outcome

	// create a new datased, children of main dataset asd
	// and populate with a template unit

	return &apiOutcome, nil
}

func (s *Server) Destroy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	var apiOutcome api.Outcome

	// destroy a solution datased, with ALL backups, snapshots

	return &apiOutcome, nil
}

func (s *Server) Backup(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Restore(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Deploy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Retreat(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Start(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Stop(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Expose(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Contain(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Snapshot(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
func (s *Server) Rollback(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log.Debugf("/ gRPS: received a message with message type: Solution")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
