package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/joegalaxy71/go-zfs"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

func (s *Server) NodeList(ctx context.Context, in *api.Void) (*api.Nodes, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodeList")

	var apinodesVal api.Nodes
	apiNodes := &apinodesVal

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		err := errors.New(message)
		return apiNodes, err
	} else {
		_log.Info("pool info obtained from onfig")
	}

	// create default master dataset name and get it via zfs wrap
	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run 'asd master init'?"
		_log.Error(message)
		_log.Error(err)
		return apiNodes, err
	} else {
		_log.Info("master dataset located")
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		return apiNodes, err
	} else {
		_log.Info("mountpoint obtained:" + mountpoint)
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db for persisting node info"
		_log.Error(message)
		_log.Error(err)
		return apiNodes, err
	} else {
		_log.Info("master db opened")
	}
	defer db.Close()

	// add node info to the k/v db
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var apiNodeVal api.Node
			apiNode := &apiNodeVal
			apiNode.Hostname = fmt.Sprint(k)
			apiNode.Ip = fmt.Sprint(v)

			apiNodes.Nodes = append(apiNodes.Nodes, apiNode)
		}
		return nil
	})
	if err != nil {
		message := "Unable to list nodes from k/v db"
		_log.Error(message)
		_log.Error(err)
		apiNodes.Outcome.Message = message
		return apiNodes, err
	}

	message := "Succesfully obtained node list"
	_log.Info(message)
	apiNodes.Outcome.Message = message
	return apiNodes, nil

}

func (s *Server) NodeAdd(ctx context.Context, in *api.Node) (*api.Node, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodeAdd")

	apiNode := in

	// a gRPC method calling another gRPC method
	conn, err := grpc.Dial(in.Ip+":9000", grpc.WithInsecure())
	if err != nil {
		message := "error dialing grpc server on asdlet on ip:" + apiNode.Ip
		log.Error(message)
		log.Error(err)
		apiNode.Outcome.Error = true
		apiNode.Outcome.Message = message
		return apiNode, err
	}
	defer conn.Close()

	c := api.NewAsdLetClient(conn)
	asdletCtx, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	apiNode, err = c.NodeAdd(asdletCtx, in)
	if err != nil {
		message := "error calling grpc:NodeAdd on asdlet on ip:" + apiNode.Ip
		log.Error(message)
		log.Error(err)
		apiNode.Outcome.Error = true
		apiNode.Outcome.Message = message
		return apiNode, err
	} else {
		// node succesfully initilized
		// proceed to add it to the k/v db
		// It will be created if it doesn't exist

		// get pool name from config
		pool := viper.GetString("pool")
		if pool == "" {
			message := "Init config value is empty for master pool"
			_log.Error(message)
			_log.Error(err)
			apiNode.Outcome.Message = message
			return apiNode, err
		}

		// create default master dataset name and get it via zfs wrap
		dataset, err := zfs.GetDataset(pool + "/asd")
		if err != nil {
			message := "Unable to locate the master dataset: did you run 'asd master init'?"
			_log.Error(message)
			_log.Error(err)
			apiNode.Outcome.Message = message
			return apiNode, err
		}

		// get the actual mountpoint
		mountpoint, err := dataset.GetProperty("mountpoint")
		if err != nil {
			message := "Unable to locate the mountpoint of the master dataset"
			_log.Error(message)
			_log.Error(err)
			apiNode.Outcome.Message = message
			return apiNode, err
		}

		// open or create the k/v db
		db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
		if err != nil {
			message := "Unable to open the master db for persisting node info"
			_log.Error(message)
			_log.Error(err)
			apiNode.Outcome.Message = message
			return apiNode, err
		}
		defer db.Close()

		// add node info to the k/v db
		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("nodes"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			db.Update(func(tx *bolt.Tx) error {
				err := b.Put([]byte(apiNode.Hostname), []byte(apiNode.Ip))
				return err
			})
			return nil
		})
		if err != nil {
			message := "Unable to update db to persist node info"
			_log.Error(message)
			_log.Error(err)
			apiNode.Outcome.Message = message
			return apiNode, err
		}

		message := "Succesfully added ADS node"
		_log.Info(message)
		apiNode.Outcome.Message = message
		return apiNode, nil

	}
}

func (s *Server) NodeRemove(ctx context.Context, in *api.Node) (*api.Node, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodeRemove")

	apiNode := in

	// proceed to remove node (by host name) it to the k/v db
	var err error

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run init?"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db to remove node"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
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
			if key == apiNode.Hostname {
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
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	message := "Succesfully removed ADS node"
	_log.Info(message)
	apiNode.Outcome.Message = message
	return apiNode, err
}

func (s *Server) NodePurge(ctx context.Context, in *api.Node) (*api.Node, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: NodePurge")

	apiNode := in

	// proceed to remove node (by host name) it to the k/v db
	var err error

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run init?"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the master db to purge node"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
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
			if key == apiNode.Hostname {
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
		apiNode.Outcome.Message = message
		return apiNode, err
	}

	message := "Succesfully purged ADS node"
	_log.Info(message)
	apiNode.Outcome.Message = message
	return apiNode, err
}
