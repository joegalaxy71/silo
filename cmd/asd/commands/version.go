package commands

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/version"
	"context"
	"github.com/op/go-logging"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var _log *logging.Logger

func Version(cmd *cobra.Command, args []string) {

	// init logs ====================================
	_log = helpers.InitLogs(true)
	_log.Debug("Version called")

	_log.Infof("ASD (client) version:%s\n", version.VERSION)

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on asdd")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewAsdClient(conn)
	var apiVoid api.Void
	apiOutcome, err := c.Version(context.Background(), &apiVoid)
	if err != nil {
		_log.Error("Unable to call ASDD gRPC server")
		_log.Error(err)
		return
	}

	_log.Infof("ASDD (server) version:%s\n", apiOutcome.Message)

	if apiOutcome.Message != version.VERSION {
		_log.Warning("Warning: client and server versions are different")
	}
}
