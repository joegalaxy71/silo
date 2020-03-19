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

	///////////////////
	// master

	var cmdMaster = &cobra.Command{
		Use:   "master",
		Short: "Manage master node",
		Long:  "Group all master subcommands",
		Args:  cobra.ExactArgs(1),
		Run:   commands.Master,
	}

	var cmdMasterInit = &cobra.Command{
		Use:   "init [pool_name]",
		Short: "Initialize the master",
		Long:  "Initialize the master by creating config file, dataset and k/v store",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MasterInit,
	}

	///////////////////////
	// solution

	var cmdSolution = &cobra.Command{
		Use:     "solution",
		Aliases: []string{"sol", "so"},
		Short:   "Subcomand for solution management",
		Long:    "Please specify the operation needed",
		Args:    cobra.ExactArgs(1),
		Run:     commands.Solution,
	}

	var cmdSolutionList = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: "list available solutions",
		Short:   "List all solutions",
		Long:    "Lists all solution with detailed info",
		Args:    cobra.NoArgs,
		Run:     commands.SolutionList,
	}

	var cmdSolutionCreate = &cobra.Command{
		Use:     "create",
		Example: "asd solution create [solution_unique_name]",
		Short:   "Create a solution",
		Long:    "Creates a solution with the given name",
		Args:    cobra.ExactArgs(1),
		Run:     commands.SolutionCreate,
	}

	var cmdSolutionCopy = &cobra.Command{
		Use:     "copy",
		Example: "asd solution copy [solution_source] [solution_dest]",
		Short:   "Copy a solution",
		Long:    "Creates a copy of a solution",
		Args:    cobra.ExactArgs(2),
		Run:     commands.SolutionCopy,
	}

	var cmdSolutionDestroy = &cobra.Command{
		Use:     "destroy",
		Example: "asd solution destroy [solution_unique_name]",
		Short:   "Destroy a solution",
		Long:    "Destroy a solution with the given name",
		Args:    cobra.ExactArgs(1),
		Run:     commands.SolutionDestroy,
	}

	var cmdSolutionDeploy = &cobra.Command{
		Use:     "deploy",
		Example: "asd solution deploy [solution_unique_name] [host_unique_name]",
		Short:   "Deploy a solution",
		Long:    "Deploy a solution with the given name in the indicated host",
		Args:    cobra.ExactArgs(1),
		Run:     commands.SolutionDeploy,
	}

	//////////////////////////////
	// node

	var cmdNode = &cobra.Command{
		Use:     "node",
		Aliases: []string{"no", "nd"},
		Short:   "Manage simple nodes",
		Long:    "Add a new node, remove and purge added ones",
		Args:    cobra.ExactArgs(1),
		Run:     commands.Node,
	}

	var cmdNodeList = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists all nodes",
		Long:    "Lists all nodes with detailed info",
		Args:    cobra.NoArgs,
		Run:     commands.NodeList,
	}

	var cmdNodeAdd = &cobra.Command{
		Use:   "add",
		Short: "Adds a node",
		Long:  "Adds a new node, initializing dataset if needed",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodeAdd,
	}

	var cmdNodeRem = &cobra.Command{
		Use:   "remove",
		Short: "Remove a node",
		Long:  "Remove an existing node, leave all solution data in place",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodeRemove,
	}

	var cmdNodePurge = &cobra.Command{
		Use:   "purge",
		Short: "Purge a node",
		Long:  "Remove an existing node, destroys all solution data",
		Args:  cobra.ExactArgs(1),
		Run:   commands.NodePurge,
	}

	var verbose bool
	var rootCmd = &cobra.Command{Use: "asd"}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(cmdVersion, cmdMaster, cmdNode, cmdSolution)
	cmdMaster.AddCommand(cmdMasterInit)
	cmdNode.AddCommand(cmdNodeList, cmdNodeAdd, cmdNodeRem, cmdNodePurge)
	cmdSolution.AddCommand(cmdSolutionList, cmdSolutionCreate, cmdSolutionCopy, cmdSolutionDestroy, cmdSolutionDeploy)
	rootCmd.Execute()
	return nil
}
