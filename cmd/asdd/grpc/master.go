package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
	"github.com/spf13/viper"
)

func (s *Server) MasterInit(ctx context.Context, in *api.Pool) (*api.Outcome, error) {
	var apiOutcome api.Outcome
	_log := helpers.InitLogs(true)
	_log.Debugf("gRPC call: Init(%s)\n", in.Name)

	_, err := zfs.CreateFilesystem(in.Name+"/asd", nil)
	if err != nil {
		message := "Error creating initial dataset " + in.Name + "/asd"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		return &apiOutcome, err
	} else {
		message := "Created root dataset " + in.Name + "/asd"
		_log.Info(message)
		apiOutcome.Error = false
		apiOutcome.Message = message
	}

	_log.Info()

	viper.Set("pool", in.Name)
	err = viper.WriteConfig()
	if err != nil {
		message := "Error persisting configuration"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		return &apiOutcome, err
	}
	return &apiOutcome, nil
}
