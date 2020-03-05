package main

import (
	"asd/common/helpers"
	_ "database/sql"
	_ "expvar" // Register the expvar handlers
	"fmt"
	"github.com/ardanlabs/conf"
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// =====================================================================================================================
//configuration
var _cfg struct {
	Web struct {
		DebugHost string `conf:"default:0.0.0.0"`
		DebugPort string `conf:"default:3300"`

		MetricsHost string `conf:"default:0.0.0.0"`
		MetricsPort string `conf:"default:2112"`

		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Log struct {
		Verbose bool `conf:"default:true"`
	}
	Server struct {
		Address string `conf:"default:asd.avero.it"`
		Port    string `conf:"default:7777"`
	}
}

// =====================================================================================================================
// version and build info

var Build int
var Version string

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
	// set version, build number
	Version = "0.9.1-303b"
	Build = 9_1_303

	// =====================================================================================================================
	// ██████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗
	//██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝
	//██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗
	//██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║
	//╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝
	//╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝

	if err := conf.Parse(os.Args[1:], "ASD-CLIENT", &_cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("ASD-CLIENT", &_cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =====================================================================================================================
	//██╗      ██████╗  ██████╗ ███████╗
	//██║     ██╔═══██╗██╔════╝ ██╔════╝
	//██║     ██║   ██║██║  ███╗███████╗
	//██║     ██║   ██║██║   ██║╚════██║
	//███████╗╚██████╔╝╚██████╔╝███████║
	//╚══════╝ ╚═════╝  ╚═════╝ ╚══════╝

	_log = helpers.InitLogs(_cfg.Log.Verbose)
	_log.Info("Log: initialized")

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

	// ██████╗ ██████╗ ██████╗  ██████╗     ██████╗██╗     ██╗███████╗███╗   ██╗████████╗
	//██╔════╝ ██╔══██╗██╔══██╗██╔════╝    ██╔════╝██║     ██║██╔════╝████╗  ██║╚══██╔══╝
	//██║  ███╗██████╔╝██████╔╝██║         ██║     ██║     ██║█████╗  ██╔██╗ ██║   ██║
	//██║   ██║██╔══██╗██╔═══╝ ██║         ██║     ██║     ██║██╔══╝  ██║╚██╗██║   ██║
	//╚██████╔╝██║  ██║██║     ╚██████╗    ╚██████╗███████╗██║███████╗██║ ╚████║   ██║
	//╚═════╝ ╚═╝  ╚═╝╚═╝      ╚═════╝     ╚═════╝╚══════╝╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝

	var err error
	conn, err = grpc.Dial(_cfg.Server.Address+":9000", grpc.WithInsecure())
	if err != nil {
		_log.Errorf("error dialing gRPC server on port 9000: %s", err)
		return err
	}

	// ██████╗ ██████╗ ███╗   ███╗███╗   ███╗ █████╗ ███╗   ██╗██████╗ ███████╗
	//██╔════╝██╔═══██╗████╗ ████║████╗ ████║██╔══██╗████╗  ██║██╔══██╗██╔════╝
	//██║     ██║   ██║██╔████╔██║██╔████╔██║███████║██╔██╗ ██║██║  ██║███████╗
	//██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██╔══██║██║╚██╗██║██║  ██║╚════██║
	//╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║  ██║██║ ╚████║██████╔╝███████║
	//╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝ ╚══════╝

	var cmdAccount = &cobra.Command{
		Use:   "account",
		Short: "Manage account on sherpa cloud",
		Long:  "Account is a master _command used to signup, signin, change or recover password and add or remove machines.",
		Args:  cobra.MinimumNArgs(0),
		Run:   account,
	}

	var cmdAccountInfo = &cobra.Command{
		Use:   "info",
		Short: "Gives back summarized account info",
		Long:  "Account info reports the number of connected machines, with summarized details about the sherpa daemons running on them.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountInfo,
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  "Prints the git commit number as build version and build date",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdVersion,
	}

	var rootCmd = &cobra.Command{Use: "sherpa"}
	rootCmd.AddCommand(cmdAccount, cmdHistory, cmdPrompt, cmdTest, cmdDaemonize, cmdVersion)
	cmdAccount.AddCommand(cmdAccountInfo, cmdAccountCreate, cmdAccountLogin)
	cmdAccount.AddCommand(cmdAccountPassword)
	cmdAccountPassword.AddCommand(cmdAccountPasswordChange, cmdAccountPasswordRecover, cmdAccountPasswordReset)
	rootCmd.Execute()

}
