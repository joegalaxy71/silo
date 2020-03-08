package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mistifyio/go-zfs"
	"github.com/spf13/viper"
	"time"
)

func (s *Server) NodeAdd(ctx context.Context, in *api.Node) (*api.NodeInfo, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodeAdd")

	var nodeInfo api.NodeInfo
	apiNodeInfo := &nodeInfo

	// node succesfully initilized
	// proceed to add it to the k/v db
	// It will be created if it doesn't exist

	var err error
	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// create default master dataset name and get it via zfs wrap
	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run init?"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db for persisting node info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}
	defer db.Close()

	// add node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		db.Update(func(tx *bolt.Tx) error {
			err := b.Put([]byte(apiNodeInfo.Hostname), []byte(apiNodeInfo.Ip))
			return err
		})
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist node info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	message := "Succesfully added ADS node"
	_log.Info(message)
	apiOutcome.Message = message
	return apiOutcome, nil

}

func (s *Server) NodeRemove(ctx context.Context, in *api.Node) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodeRemove")

	var outcome api.Outcome
	apiOutcome := &outcome

	// proceed to remove node (by host name) it to the k/v db
	var err error

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run init?"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db to remove node"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}
	defer db.Close()

	// remove node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		c := b.Cursor()
		deleted := false

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := fmt.Sprint(k)
			if key == in.Host {
				c.Delete()
				deleted = true
			}
			if deleted == false {
				message := "node not found"
				_log.Error(message)
				err := errors.New(message)
				return err
			}
			//fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist node info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	message := "Succesfully removed ADS node"
	_log.Info(message)
	apiOutcome.Message = message
	return apiOutcome, nil
}

func (s *Server) NodePurge(ctx context.Context, in *api.Node) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodePurge")

	var outcome api.Outcome
	apiOutcome := &outcome

	// proceed to remove node (by host name) it to the k/v db
	var err error

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run init?"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db to remove node"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}
	defer db.Close()

	// remove node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		c := b.Cursor()
		deleted := false

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := fmt.Sprint(k)
			if key == in.Host {
				c.Delete()
				deleted = true
			}
			if deleted == false {
				message := "node to purge not found"
				_log.Error(message)
				err := errors.New(message)
				return err
			}
			//fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist node info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	}

	message := "Succesfully removed ADS node"
	_log.Info(message)
	apiOutcome.Message = message
	return apiOutcome, nil
}
