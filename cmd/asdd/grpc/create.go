package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
)

func (s *Server) Create(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")
	var apiOutcome api.Outcome

	// create a new datased, children of main dataset asd
	// and populate with a template unit

	return &apiOutcome, nil
}
