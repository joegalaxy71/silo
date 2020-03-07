package main

import (
	"asd/cmd/asd/commands"
	"asd/common/helpers"
	_ "database/sql"
	_ "expvar" // Register the expvar handlers
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

// =====================================================================================================================
// module wide globals (logging, locks, db, etc)
var _log *logging.Logger
var _lock sync.Mutex

// =====================================================================================================================
//███╗   ███╗ █████╗ ██╗███╗   ██╗
//████╗ ████║██╔══██╗██║████╗  ██║
//██╔████╔██║███████║██║██╔██╗ ██║
//██║╚██╔╝██║██╔══██║██║██║╚██╗██║
//██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
//╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝
func main() {
	if err := run(); err != nil {
		log.Println("run() function error: ", err)
	}
}

// =====================================================================================================================
//██████╗ ██╗   ██╗███╗   ██╗
//██╔══██╗██║   ██║████╗  ██║
//██████╔╝██║   ██║██╔██╗ ██║
//██╔══██╗██║   ██║██║╚██╗██║
//██║  ██║╚██████╔╝██║ ╚████║
//╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
func run() error {

	// =====================================================================================================================
	//██╗      ██████╗  ██████╗ ███████╗
	//██║     ██╔═══██╗██╔════╝ ██╔════╝
	//██║     ██║   ██║██║  ███╗███████╗
	//██║     ██║   ██║██║   ██║╚════██║
	//███████╗╚██████╔╝╚██████╔╝███████║
	//╚══════╝ ╚═════╝  ╚═════╝ ╚══════╝

	_log = helpers.InitLogs(true)

	// =====================================================================================================================
	//███████╗██╗ ██████╗ ███╗   ██╗ █████╗ ██╗     ███████╗
	//██╔════╝██║██╔════╝ ████╗  ██║██╔══██╗██║     ██╔════╝
	//███████╗██║██║  ███╗██╔██╗ ██║███████║██║     ███████╗
	//╚════██║██║██║   ██║██║╚██╗██║██╔══██║██║     ╚════██║
	//███████║██║╚██████╔╝██║ ╚████║██║  ██║███████╗███████║
	//╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚══════╝

	helpers.InitSignals()

	// =====================================================================================================================
	//██╗   ██╗███╗   ██╗██╗ ██████╗ ██╗   ██╗███████╗
	//██║   ██║████╗  ██║██║██╔═══██╗██║   ██║██╔════╝
	//██║   ██║██╔██╗ ██║██║██║   ██║██║   ██║█████╗
	//██║   ██║██║╚██╗██║██║██║▄▄ ██║██║   ██║██╔══╝
	//╚██████╔╝██║ ╚████║██║╚██████╔╝╚██████╔╝███████╗
	//╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚══▀▀═╝  ╚═════╝ ╚══════╝

	s := single.New("asd")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		log.Println("another instance of the app is already running, exiting")
		return err
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		log.Println("failed to acquire exclusive app lock: %v", err)
		return err
	}
	defer s.TryUnlock()

	// ██████╗ ██████╗ ███╗   ███╗███╗   ███╗ █████╗ ███╗   ██╗██████╗ ███████╗
	//██╔════╝██╔═══██╗████╗ ████║████╗ ████║██╔══██╗████╗  ██║██╔══██╗██╔════╝
	//██║     ██║   ██║██╔████╔██║██╔████╔██║███████║██╔██╗ ██║██║  ██║███████╗
	//██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██╔══██║██║╚██╗██║██║  ██║╚════██║
	//╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║  ██║██║ ╚████║██████╔╝███████║
	//╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝ ╚══════╝

	//var cmdAccount = &cobra.Command{
	//	Use:   "account",
	//	Short: "Manage account on sherpa cloud",
	//	Long:  "Account is a master _command used to signup, signin, change or recover password and add or remove machines.",
	//	Args:  cobra.MinimumNArgs(0),
	//	Run:   account,
	//}
	//
	//var cmdAccountInfo = &cobra.Command{
	//	Use:   "info",
	//	Short: "Gives back summarized account info",
	//	Long:  "Account info reports the number of connected machines, with summarized details about the sherpa daemons running on them.",
	//	Args:  cobra.MinimumNArgs(0),
	//	Run:   accountInfo,
	//}

	var cmdVersion = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Prints version information",
		Long:    "Prints the git commit number as build version and build date",
		Args:    cobra.MinimumNArgs(0),
		Run:     commands.Version,
	}

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialize ASD",
		Long:  "Creates necessary datasets and config for ASD",
		Args:  cobra.ExactArgs(1),
		Run:   commands.Init,
	}

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a solution",
		Long:  "Creates a solution with the given name",
		Args:  cobra.ExactArgs(1),
		Run:   commands.Create,
	}

	var cmdNode = &cobra.Command{
		Use:   "node",
		Short: "Manage adding and removing nodes",
		Long:  "Add a new node, remove and purge added ones",
		Args:  cobra.MinimumNArgs(0),
		Run:   commands.Node,
	}

	var cmdNodeAdd = &cobra.Command{
		Use:   "node add [node_ip]",
		Short: "Add an ADS node",
		Long:  "Add a new node, initalialyzing datasets if needed",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodeAdd,
	}

	var cmdNodeRem = &cobra.Command{
		Use:   "node remove [node_ip]",
		Short: "Remove an ADS node",
		Long:  "Remove an existing node, leave all solution data in place",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodeRemove,
	}

	var cmdNodePurge = &cobra.Command{
		Use:   "node purge [node_ip]",
		Short: "Purge an ADS node",
		Long:  "Remove an existing node, destroys all solution data",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodePurge,
	}

	var verbose bool
	var rootCmd = &cobra.Command{Use: "asd"}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(cmdVersion, cmdInit, cmdCreate, cmdNode)
	cmdNode.AddCommand(cmdNodeAdd, cmdNodeRem, cmdNodePurge)
	//cmdAccount.AddCommand(cmdAccountInfo, cmdAccountCreate, cmdAccountLogin)
	//cmdAccount.AddCommand(cmdAccountPassword)
	//cmdAccountPassword.AddCommand(cmdAccountPasswordChange, cmdAccountPasswordRecover, cmdAccountPasswordReset)
	rootCmd.Execute()

	return nil
}
