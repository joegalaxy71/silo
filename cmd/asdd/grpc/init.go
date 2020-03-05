package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
)

func (s *Server) Init(ctx context.Context, in *api.Pool) (*api.Outcome, error) {
	var apiOutcome api.Outcome
	_log := helpers.InitLogs(true)
	_log.Debugf("gRPC call: Init(%s)\n", in.Name)

	_, err := zfs.CreateFilesystem(in.Name+"/asd", nil)
	if err != nil {
		_log.Error("Error creating initial dataset " + in.Name + "/asd")
		apiOutcome.Error = true
		apiOutcome.Message = "Error creating initial dataset " + in.Name + "/asd"
		return &apiOutcome, err
	}

	println("Created zfs volume" + in.Name + "/asd")

	return &apiOutcome, nil
}
