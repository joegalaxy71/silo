package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NodeAdd(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:NodeAdd")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()

	c := api.NewAsddClient(conn)
	var apiSolution api.Solution
	apiSolution.Name = args[0]
	apiOutcome, err := c.Create(context.Background(), &apiSolution)
	if err != nil {
		_log.Error("Unable to add the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiOutcome.Message)
}

func NodeRemove(cmd *cobra.Command, args []string) {
	_log = helpers.InitLogs(true)
	_log.Debug("Command:NodeRemove")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()

	c := api.NewAsddClient(conn)
	var apiNode api.Node
	apiNode.Ip = args[0]
	apiOutcome, err := c.NodeRemove(context.Background(), &apiNode)
	if err != nil {
		_log.Error("Unable to remove the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiOutcome.Message)
}

func NodePurge(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:NodePurge")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()

	c := api.NewAsddClient(conn)
	var apiNode api.Node
	apiNode.Ip = args[0]
	apiOutcome, err := c.NodePurge(context.Background(), &apiNode)
	if err != nil {
		_log.Error("Unable to purge the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiOutcome.Message)
}
