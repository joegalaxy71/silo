package grpc

import (
	"asd/common/api"
	"context"
	"errors"
)

func (s *Server) Version(ctx context.Context, in *api.Void) (*api.Outcome, error) {

	var apiOutcomeVal api.Outcome
	apiOutcome := &apiOutcomeVal

	apiOutcome.Error = true
	apiOutcome.Message = "Not implemented"

	err := errors.New(apiOutcome.Message)

	return apiOutcome, err
}
