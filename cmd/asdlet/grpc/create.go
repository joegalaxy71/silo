package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
	"github.com/spf13/viper"
)

func (s *Server) Create(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")

	var pool string
	pool = viper.GetString("pool")
	var apiOutcome api.Outcome

	datasetName := pool + "/asd/" + in.Name

	_, err := zfs.CreateFilesystem(datasetName, nil)
	if err != nil {
		_log.Error("Error creating solution " + datasetName)
		apiOutcome.Error = true
		apiOutcome.Message = "Error creating solution " + datasetName
		return &apiOutcome, err
	} else {
		apiOutcome.Error = false
		apiOutcome.Message = "Succesfully created solution " + datasetName
		return &apiOutcome, nil
	}

}
