package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/zfs"
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

func (s *Server) SolutionList(ctx context.Context, in *api.Void) (*api.Solutions, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: SolutionList")

	var apiSolutions *api.Solutions
	var apiOutcome api.Outcome
	apiSolutions.Outcome = &apiOutcome

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		err := errors.New(message)
		return apiSolutions, err
	} else {
		_log.Info("got pool from config file:" + pool)
	}

	// create default master dataset name and get it via zfs wrap
	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run 'asd master init'?"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("Master dataset found:" + dataset.Name)
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("Got mountpoint: " + mountpoint)
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting node info"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("succesfully opened main db")
	}
	defer db.Close()

	// add node info to the k/v db
	err = db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			var apiSolutionVal api.Solution
			apiSolution := &apiSolutionVal
			apiSolution.Name = fmt.Sprint(k)

			apiSolutions.Solutions = append(apiSolutions.Solutions, apiSolution)
		}
		return nil
	})
	if err != nil {
		message := "Unable to list solutions from k/v db"
		_log.Error(message)
		_log.Error(err)
		apiSolutions.Outcome.Message = message
		return apiSolutions, err
	}

	message := "Succesfully obtained solution list"
	_log.Info(message)
	apiSolutions.Outcome.Message = message
	return apiSolutions, nil
}

func (s *Server) SolutionCreate(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")

	apiSolution := in
	var apiOutcome api.Outcome
	apiSolution.Outcome = &apiOutcome

	var pool string
	pool = viper.GetString("pool")
	if pool == "" {
		message := "The master pool is unconfigured"
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	datasetName := pool + "/asd/" + apiSolution.Name

	dataset, err := zfs.CreateFilesystem(datasetName, nil)
	if err != nil {
		message := "Error creating dataset for new solutions" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Succesfully created dataset:" + datasetName)
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Got mountpoint:" + mountpoint)
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db opened succesfully")
	}

	defer db.Close()

	// add node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte(apiSolution.Name), []byte(""))
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}
		return nil
	})
	if err != nil {
		message := "Unable to update db to persist solution info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db updated with master info")
	}

	message := "Succesfully created new solution" + apiSolution.Name
	_log.Info(message)
	apiSolution.Outcome.Error = false
	apiSolution.Outcome.Message = message
	return apiSolution, nil
}

func (s *Server) SolutionDestroy(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Destroy")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionDeploy(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Deploy")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionRetire(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Retreat")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionStart(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Start")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionStop(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Stop")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionSnapshot(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Snapshot")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionRollback(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Rollback")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionBackup(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Backup")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionRestore(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Restore")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionExpose(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Expose")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionContain(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Contain")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}
