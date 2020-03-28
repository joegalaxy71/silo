package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"context"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"text/tabwriter"
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
	//_log.Infof("Solution, host, status")
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "NAME\tHOSTNAME\tSTATUS\t")
	for _, apiSolution := range apiSolutions.Solutions {
		fmt.Fprintln(w, apiSolution.Name+"\t"+apiSolution.Hostname+"\t"+apiSolution.Status+"\t")
		//_log.Infof("%s, %s, %s\n", apiSolution.Name, apiSolution.Hostname, apiSolution.Status)
	}
	w.Flush()
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

func SolutionCopy(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:SolutionCopy")

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewAsddClient(conn)
	var apiCopyArgsVal api.CopyArgs
	apiCopyArgs := &apiCopyArgsVal
	apiCopyArgs.Source = args[0]
	apiCopyArgs.Destination = args[1]
	apiOutcome, err := c.SolutionCopy(context.Background(), apiCopyArgs)
	if err != nil {
		_log.Error("Adding solution failed")
		_log.Error(err)
		return
	}

	_log.Infof("Outcome message:%s\n", apiOutcome.Message)
}

func SolutionDestroy(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:SolutionDestroy")

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
	apiSolution, err = c.SolutionDestroy(context.Background(), apiSolution)
	if err != nil {
		_log.Error("Destroying solution failed")
		_log.Error(err)
		return
	}

	_log.Infof("Outcome message:%s\n", apiSolution.Outcome.Message)
}

func SolutionDeploy(cmd *cobra.Command, args []string) {

	_log = helpers.InitLogs(true)
	_log.Debug("Command:SolutionDeploy")

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
	apiSolution.Hostname = args[1]
	apiSolution, err = c.SolutionDeploy(context.Background(), apiSolution)
	if err != nil {
		_log.Error("Deploy failed")
		_log.Error(err)
		return
	}

	_log.Infof("Outcome message:%s\n", apiSolution.Outcome.Message)
}
