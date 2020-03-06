package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Create(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:Create")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewAsdClient(conn)
	var apiSolution api.Solution
	apiSolution.Name = args[0]
	apiOutcome, err := c.Create(context.Background(), &apiSolution)
	if err != nil {
		_log.Error("Unable to call ASDD gRPC server")
		_log.Error(err)
		return
	}

	_log.Infof("Outcome message:%s\n", apiOutcome.Message)
}
