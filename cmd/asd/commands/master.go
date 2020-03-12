package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Master(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:Master")

	_log.Error("Please call 'master' with more parameters")
}

func MasterInit(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewAsddClient(conn)

	var apiMasterVal api.Master
	apiMaster := &apiMasterVal

	apiMaster.Poolname = args[0]
	apiMaster, err = c.MasterInit(context.Background(), apiMaster)
	if err != nil {
		_log.Error("master init command failed")
		_log.Error(err)
		return
	} else {
		_log.Infof(apiMaster.Outcome.Message)
	}
}
