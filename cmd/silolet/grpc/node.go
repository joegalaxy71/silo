package grpc

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"os"
	"silo/common/api"
	"silo/common/helpers"
	"silo/common/zfs"
)

func (s *Server) NodeAdd(ctx context.Context, in *api.Node) (*api.Node, error) {

	apiNode := in

	_log := helpers.InitLogs(true)
	_log.Debugf("gRPC call: NodeAdd(pool:%s)\n", apiNode.Poolname)

	_, err := zfs.CreateFilesystem(in.Poolname+"/silo", nil)
	if err != nil {
		message := "Error creating initial dataset " + apiNode.Poolname + "/silo"
		_log.Error(message)
		apiNode.Outcome.Error = true
		apiNode.Outcome.Message = message
		return apiNode, err
	} else {
		message := "Created root dataset " + apiNode.Poolname + "/asd"
		_log.Info(message)
		apiNode.Outcome.Error = false
		apiNode.Outcome.Message = message
	}

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		message := "Unable to get master hostname"
		_log.Error(message)
		_log.Error(err)
		apiNode.Outcome.Message = message
		return apiNode, err
	} else {
		_log.Info("Got master hostname:" + hostname)
		apiNode.Hostname = hostname
	}

	viper.Set("pool", apiNode.Poolname)
	err = viper.WriteConfig()
	if err != nil {
		message := "Error persisting configuration"
		_log.Error(message)
		apiNode.Outcome.Error = true
		apiNode.Outcome.Message = message
		return apiNode, err
	}
	return apiNode, nil
}

func (s *Server) NodeRemove(ctx context.Context, in *api.Node) (*api.Node, error) {

	apiNode := in

	apiNode.Outcome.Error = true
	apiNode.Outcome.Message = "Not implemented"

	err := errors.New(apiNode.Outcome.Message)

	return apiNode, err
}

func (s *Server) NodePurge(ctx context.Context, in *api.Node) (*api.Node, error) {

	apiNode := in

	apiNode.Outcome.Error = true
	apiNode.Outcome.Message = "Not implemented"

	err := errors.New(apiNode.Outcome.Message)

	return apiNode, err
}
