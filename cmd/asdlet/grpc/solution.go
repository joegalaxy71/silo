package grpc

import (
	"asd/common/api"
	"context"
	"errors"
)

func (s *Server) SolutionDeploy(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}

func (s *Server) SolutionRetire(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}

func (s *Server) SolutionStart(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}

func (s *Server) SolutionStop(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}

func (s *Server) SolutionExpose(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}

func (s *Server) SolutionContain(ctx context.Context, in *api.Solution) (*api.Solution, error) {

	apiSolution := in

	apiSolution.Outcome.Error = true
	apiSolution.Outcome.Message = "Not implemented"

	err := errors.New(apiSolution.Outcome.Message)

	return apiSolution, err
}
