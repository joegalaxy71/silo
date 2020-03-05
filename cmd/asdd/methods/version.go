package methods

import (
	"asd/common/api"
	"asd/common/version"
	"context"
)

func (s *Server) Version(ctx context.Context, in *api.Void) (*api.Outcome, error) {
	_log.Debugf("method VERSION called")

	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome
	apiOutcome.Error = false
	apiOutcome.Message = version.VERSION

	return &apiOutcome, nil
}
