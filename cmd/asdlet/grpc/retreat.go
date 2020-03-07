package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
)

func (s *Server) Retreat(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Retreat")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
