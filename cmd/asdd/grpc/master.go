package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
	"github.com/spf13/viper"
)

func (s *Server) MasterInit(ctx context.Context, in *api.Master) (*api.Master, error) {
	apiMaster := in
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: MasterInit(%s)\n")

	_, err := zfs.CreateFilesystem(apiMaster.Poolname+"/asd", nil)
	if err != nil {
		message := "Error creating initial dataset " + apiMaster.Poolname + "/asd"
		_log.Error(message)
		apiMaster.Outcome.Error = true
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		message := "Created root dataset " + apiMaster.Poolname + "/asd"
		_log.Info(message)
		apiMaster.Outcome.Error = false
		apiMaster.Outcome.Message = message
	}

	_log.Info()

	viper.Set("pool", apiMaster.Poolname)
	err = viper.WriteConfig()
	if err != nil {
		message := "Error persisting configuration"
		_log.Error(message)
		apiMaster.Outcome.Error = true
		apiMaster.Outcome.Message = message
		return apiMaster, err
	}
	return apiMaster, nil
}
