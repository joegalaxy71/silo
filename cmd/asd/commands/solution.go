package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Solution(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:Solution")

	_log.Error("Please call 'solution' with more parameters")
}

func SolutionList(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:SolutionList")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()

	c := api.NewAsddClient(conn)
	var apiVoidVal api.Void
	apiVoid := &apiVoidVal
	apiSolutions, err := c.SolutionList(context.Background(), apiVoid)
	if err != nil {
		_log.Error("Unable to list available solutions, detailed error message follows")
		_log.Error(err)
		return
	}

	_log.Info(apiSolutions.Outcome.Message)
}

func SolutionCreate(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:SolutionCreate")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewAsddClient(conn)
	var apiSolutionVal api.Solution
	apiSolution := &apiSolutionVal
	apiSolution.Name = args[0]
	apiSolution, err = c.SolutionCreate(context.Background(), apiSolution)
	if err != nil {
		_log.Error("Adding solution failed")
		_log.Error(err)
		return
	}

	_log.Infof("Outcome message:%s\n", apiSolution.Outcome.Message)
}
