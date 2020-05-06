package grpc

import (
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"os"
	"silo/common/api"
	"silo/common/helpers"
	"silo/common/zfs"
	"time"
)

func (s *Server) MasterInit(ctx context.Context, in *api.Master) (*api.Master, error) {
	apiMaster := in
	var apiOutcomeVal api.Outcome
	apiOutcome := &apiOutcomeVal
	apiMaster.Outcome = apiOutcome
	_log := helpers.InitLogs(true)
	_log.Debugf("gRPC call: MasterInit(%s)\n", apiMaster.Poolname)

	dataset, err := zfs.CreateFilesystem(apiMaster.Poolname+"/silo", nil)
	//.CreateFilesystem(apiMaster.Poolname+"/silo", nil)
	if err != nil {
		message := "Error creating initial dataset " + apiMaster.Poolname + "/silo"
		_log.Error(message)
		apiMaster.Outcome.Error = true
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		_log.Info("Created root dataset " + apiMaster.Poolname + "/silo")
	}

	//get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		return apiMaster, err
	} else {
		_log.Info("Got mountpoint property:  " + mountpoint)
	}

	viper.Set("mountpoint", mountpoint)
	err = viper.WriteConfig()
	if err != nil {
		message := "Error persisting configuration (mountpoint)"
		_log.Error(message)
		apiMaster.Outcome.Error = true
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		_log.Info("Mountpoint property updated on " + viper.ConfigFileUsed())
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
		_log.Info("Pool property updated on " + viper.ConfigFileUsed())
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/silo.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		_log.Info("Main db opened succesfully")
	}

	defer db.Close()

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		message := "Unable to get master hostname"
		_log.Error(message)
		_log.Error(err)
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		_log.Info("Got master hostname:" + hostname)
		apiMaster.Hostname = hostname
	}

	// add node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("maste"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		encoded, err := proto.Marshal(apiMaster)
		if err != nil {
			return err
		}

		err = b.Put([]byte(apiMaster.Hostname), encoded)
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist master info"
		_log.Error(message)
		_log.Error(err)
		apiMaster.Outcome.Message = message
		return apiMaster, err
	} else {
		_log.Info("Main db updated with master info")
	}

	message := "Succesfully initialized silo master"
	_log.Info(message)
	apiMaster.Outcome.Message = message
	return apiMaster, nil
}
