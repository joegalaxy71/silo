package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/zfs"
	"context"
	"fmt"
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
	} else {
		message := "Got mountpoint property:  " + mountpoint
		_log.Info(message)
		apiMaster.Outcome.Error = false
		apiMaster.Outcome.Message = message
	}

	viper.Set("pool", apiMaster.Poolname)
	err = viper.WriteConfig()
	if err != nil {
		message := "Error persisting configuration"
		_log.Error(message)
		apiMaster.Outcome.Error = true
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		message := "Configuration updated"
		_log.Info(message)
		apiMaster.Outcome.Error = false
		apiMaster.Outcome.Message = message
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		message := "Main db opened succesfully"
		_log.Info(message)
		apiMaster.Outcome.Error = false
		apiMaster.Outcome.Message = message
	}

	defer db.Close()

	// add node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("masters"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		db.Update(func(tx *bolt.Tx) error {
			err := b.Put([]byte(apiMaster.Hostname), []byte(apiMaster.Ip))
			return err
		})
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist master info"
		_log.Error(message)
		_log.Error(err)
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		message := "Main db updated with master info"
		_log.Info(message)
		_log.Info(err)
		apiMaster.Outcome.Message = message
	}

	message := "Succesfully initialized ADS master"
	_log.Info(message)
	apiMaster.Outcome.Message = message
	return apiMaster, nil
}
