package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/zfs"
	"context"
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
	"time"
)

func (s *Server) MasterInit(ctx context.Context, in *api.Master) (*api.Master, error) {
	apiMaster := in
	var apiOutcomeVal api.Outcome
	apiOutcome := &apiOutcomeVal
	apiMaster.Outcome = apiOutcome
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: MasterInit(%s)\n")

	dataset, err := zfs.CreateFilesystem(apiMaster.Poolname+"/asd", nil)
	//.CreateFilesystem(apiMaster.Poolname+"/asd", nil)
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

	//get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		return apiMaster, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 20 * time.Second})
	if err != nil {
		message := "Unable to open/create the master db for persisting node info"
		_log.Error(message)
		_log.Error(err)
		return apiMaster, err
	}
	defer db.Close()

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
