package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
)

func (s *Server) Init(ctx context.Context, in *api.Pool) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Init")
	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
