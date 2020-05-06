package grpc

import (
	"context"
	"silo/common/api"
	"silo/common/helpers"
	"silo/common/version"
)

func (s *Server) Version(ctx context.Context, in *api.Void) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Version")

	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome
	apiOutcome.Error = false
	apiOutcome.Message = version.VERSION

	return &apiOutcome, nil
}
