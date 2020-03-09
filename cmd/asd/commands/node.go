package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Node(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:Node")

	_log.Error("Please call 'node' with more parameters")
}

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
	var apiNodeVal api.Node
	apiNode := &apiNodeVal
	apiNode.Ip = args[0]
	apiNode, err = c.NodeAdd(context.Background(), apiNode)
	if err != nil {
		_log.Error("Unable to add the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiNode.Outcome.Message)
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
	var apiNodeVal api.Node
	apiNode := &apiNodeVal
	apiNode.Ip = args[0]
	apiNode, err = c.NodeRemove(context.Background(), apiNode)
	if err != nil {
		_log.Error("Unable to remove the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiNode.Outcome.Message)
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
	var apiNodeVal api.Node
	apiNode := &apiNodeVal
	apiNode.Ip = args[0]
	apiNode, err = c.NodePurge(context.Background(), apiNode)
	if err != nil {
		_log.Error("Unable to purge the specified node, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiNode.Outcome.Message)
}