package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
)

func (s *Server) Create(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")
	var apiOutcome api.Outcome

	_, err := zfs.CreateFilesystem(in.Name+"/asd", nil)
	if err != nil {
		_log.Error("Error creating initial dataset " + in.Name + "/asd")
		apiOutcome.Error = true
		apiOutcome.Message = "Error creating initial dataset " + in.Name + "/asd"
		return &apiOutcome, err
	}

	return &apiOutcome, nil
}
