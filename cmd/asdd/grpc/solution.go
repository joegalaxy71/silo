package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func (s *Server) SolutionCreate(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")

	apiSolution := in

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

	_, err := zfs.CreateFilesystem(datasetName, nil)
	if err != nil {
		message := "Error creating dataset for new solutions" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		message := "Succesfully created new solutions with dataset:" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = false
		apiSolution.Outcome.Message = message
		return apiSolution, nil
	}
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
