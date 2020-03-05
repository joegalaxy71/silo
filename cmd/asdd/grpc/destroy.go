package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
)

func (s *Server) Destroy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Destroy")

	var apiOutcome api.Outcome

	// destroy a solution datased, with ALL backups, snapshots

	return &apiOutcome, nil
}
