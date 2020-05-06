package commands

import (
	"context"
	"github.com/op/go-logging"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"silo/common/api"
	"silo/common/helpers"
	"silo/common/version"
)

var _log *logging.Logger

func Version(cmd *cobra.Command, args []string) {

	// init logs ====================================
	_log = helpers.InitLogs(true)
	_log.Debug("Command: version")

	_log.Infof("silo (client) version:%s\n", version.VERSION)

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		log.Error("error dialing grpc server on silod")
		log.Error(err)
		return
	}
	defer conn.Close()
	c := api.NewSilodClient(conn)
	var apiVoid api.Void
	apiOutcome, err := c.Version(context.Background(), &apiVoid)
	if err != nil {
		_log.Error("Unable to call siloD gRPC server")
		_log.Error(err)
		return
	}

	_log.Infof("siloD (server) version:%s\n", apiOutcome.Message)

	if apiOutcome.Message != version.VERSION {
		_log.Warning("Warning: client and server versions are different")
	}
}
