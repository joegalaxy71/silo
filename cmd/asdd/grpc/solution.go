package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/version"
	"context"
	"github.com/mistifyio/go-zfs"
	"github.com/spf13/viper"
)

func (s *Server) Version(ctx context.Context, in *api.Void) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Version")

	//should
	// create the config
	// create the main dataset: asd

	var apiOutcome api.Outcome
	apiOutcome.Error = false
	apiOutcome.Message = version.VERSION

	return &apiOutcome, nil
}

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

func (s *Server) Destroy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Destroy")

	var apiOutcome api.Outcome

	// destroy a solution datased, with ALL backups, snapshots

	return &apiOutcome, nil
}

func (s *Server) Deploy(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Deploy")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Retire(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Retreat")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Start(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Start")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Stop(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Stop")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Snapshot(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Snapshot")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Rollback(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Rollback")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Backup(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Backup")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Restore(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Restore")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Expose(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Expose")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}

func (s *Server) Contain(ctx context.Context, in *api.Solution) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Contain")

	// create var to build up api response
	var apiOutcome api.Outcome

	return &apiOutcome, nil
}
