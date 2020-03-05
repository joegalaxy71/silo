package commands

import (
	"asd/common/helpers"
	"github.com/spf13/cobra"
)

func Version(cmd *cobra.Command, args []string) {

	// init logs ====================================
	log := helpers.InitLogs(true)
	log.Info("Version called")

	//client, err := dialGrpc()
	//if err != nil {
	//	log.Error("error dialing grpc server on asdd")
	//	log.Error(err)
	//}

	//var err = errors.New("not implemented")
	//return err
}
