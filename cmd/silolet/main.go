package main

import (
	_ "expvar"
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	siloletGrpc "silo/cmd/silolet/grpc"
	"silo/common/helpers"
)

// =====================================================================================================================
// version and build info

var Version string

// =====================================================================================================================
// Server represents the gRPC server

type Server struct {
}

// =====================================================================================================================
// module wide globals (logging, locks, db, etc)
var _log *logging.Logger
var _verbose bool

// =====================================================================================================================
// metrics globals
var (
	_workerOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "silo_worker_ops_total",
		Help: "The total number worker was spawned by internal cron subsys",
	})
)

const CONFIGDIR = "/etc/silolet/"
const CONFIGFILE = "/etc/silolet/config.yaml"

// =====================================================================================================================
// main
func main() {

	if err := run(); err != nil {
		log.Println("run() function error: ", err)
	}
}

// =====================================================================================================================
// run
func run() error {

	// =====================================================================================================================
	// set version, build number
	Version = "0.3.0-3a"

	//██╗      ██████╗  ██████╗ ███████╗
	//██║     ██╔═══██╗██╔════╝ ██╔════╝
	//██║     ██║   ██║██║  ███╗███████╗
	//██║     ██║   ██║██║   ██║╚════██║
	//███████╗╚██████╔╝╚██████╔╝███████║
	//╚══════╝ ╚═════╝  ╚═════╝ ╚══════╝

	_log = helpers.InitLogs(true)
	_log.Debug("Log: initialized")

	//██████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗
	//██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝
	//██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗
	//██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║
	//╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝
	//╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝

	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(CONFIGDIR) // path to look for the config file in

	// get CONFIGDIR info
	_, err := os.Stat(CONFIGDIR)

	// if does't exists
	if os.IsNotExist(err) {
		// create it
		errDir := os.MkdirAll(CONFIGDIR, 0755)
		// if create error
		if errDir != nil {
			_log.Errorf("unable to create config directory %s\n", CONFIGDIR)
			return err
		}
	}

	// get CONFIGDIR info
	_, err = os.Stat(CONFIGFILE)

	// if does't exists
	if os.IsNotExist(err) {
		// create it
		_, errFile := os.Create(CONFIGFILE)
		// if file create error
		if errFile != nil {
			_log.Errorf("unable to create config file %s\n", CONFIGFILE)
			return err
		} else {
			//file created, setting defaults
			// correctly created dir
			defaults := map[string]interface{}{
				"pool": "undefined",
				"auth": map[string]string{
					"username": "none",
					"password": "unset",
				},
			}

			for key, value := range defaults {
				viper.SetDefault(key, value)
			}

			viper.WriteConfig()
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				_log.Error("config file " + CONFIGFILE + " not found")
				return err
			} else {
				// Config file was found but another error was produced
				_log.Error("error reading config file " + CONFIGFILE)
				return err
			}
		}
	}

	//███████╗██╗ ██████╗ ███╗   ██╗ █████╗ ██╗     ███████╗
	//██╔════╝██║██╔════╝ ████╗  ██║██╔══██╗██║     ██╔════╝
	//███████╗██║██║  ███╗██╔██╗ ██║███████║██║     ███████╗
	//╚════██║██║██║   ██║██║╚██╗██║██╔══██║██║     ╚════██║
	//███████║██║╚██████╔╝██║ ╚████║██║  ██║███████╗███████║
	//╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚══════╝
	helpers.InitSignals()
	_log.Debugf("Signals: initialized")

	//██╗   ██╗███╗   ██╗██╗ ██████╗ ██╗   ██╗███████╗
	//██║   ██║████╗  ██║██║██╔═══██╗██║   ██║██╔════╝
	//██║   ██║██╔██╗ ██║██║██║   ██║██║   ██║█████╗
	//██║   ██║██║╚██╗██║██║██║▄▄ ██║██║   ██║██╔══╝
	//╚██████╔╝██║ ╚████║██║╚██████╔╝╚██████╔╝███████╗
	//╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚══▀▀═╝  ╚═════╝ ╚══════╝

	s := single.New("asdlet")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		_log.Errorf("another instance of the app is already running, exiting")
		return err
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		_log.Errorf("failed to acquire exclusive app lock: %v", err)
		return err
	}
	defer s.TryUnlock()
	_log.Debugf("Unique: initialized")

	//██████╗ ███████╗██████╗ ██╗   ██╗ ██████╗
	//██╔══██╗██╔════╝██╔══██╗██║   ██║██╔════╝
	//██║  ██║█████╗  ██████╔╝██║   ██║██║  ███╗
	//██║  ██║██╔══╝  ██╔══██╗██║   ██║██║   ██║
	//██████╔╝███████╗██████╔╝╚██████╔╝╚██████╔╝
	//╚═════╝ ╚══════╝╚═════╝  ╚═════╝  ╚═════╝

	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.

	go func() {
		_log.Debugf("Debug: listening on: %s", "0.0.0.0:3300")
		_log.Debugf("Debug: listener closed : %v", http.ListenAndServe("0.0.0.0:3300", http.DefaultServeMux))
	}()
	_log.Debugf("Debug: initialized")

	//███╗   ███╗███████╗████████╗██████╗ ██╗ ██████╗███████╗
	//████╗ ████║██╔════╝╚══██╔══╝██╔══██╗██║██╔════╝██╔════╝
	//██╔████╔██║█████╗     ██║   ██████╔╝██║██║     ███████╗
	//██║╚██╔╝██║██╔══╝     ██║   ██╔══██╗██║██║     ╚════██║
	//██║ ╚═╝ ██║███████╗   ██║   ██║  ██║██║╚██████╗███████║
	//╚═╝     ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝ ╚═════╝╚══════╝

	// /metrics - Added to the metrics handler

	go func() {
		_log.Debugf("Metrics: listening on %s", "0.0.0.0:2112")
		http.Handle("/metrics", promhttp.Handler())
		_log.Debugf("Metrics: listener closed : %v", http.ListenAndServe("0.0.0.0:2112", nil))
	}()
	_log.Debugf("Metrics: initialized")

	// ██████╗██╗  ██╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██╗     ███████╗
	//██╔════╝██║  ██║██╔══██╗████╗  ██║████╗  ██║██╔════╝██║     ██╔════╝
	//██║     ███████║███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██║     ███████╗
	//██║     ██╔══██║██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██║     ╚════██║
	//╚██████╗██║  ██║██║  ██║██║ ╚████║██║ ╚████║███████╗███████╗███████║
	//╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝

	_log.Debugf("Channels: initialized")

	//██████╗██████╗  ██████╗ ███╗   ██╗
	//██╔════╝██╔══██╗██╔═══██╗████╗  ██║
	//██║     ██████╔╝██║   ██║██╔██╗ ██║
	//██║     ██╔══██╗██║   ██║██║╚██╗██║
	//╚██████╗██║  ██║╚██████╔╝██║ ╚████║
	//╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

	cronTab := cron.New()
	err = cronTab.AddFunc("0 * * * * *", worker) //every minute
	if err != nil {
		errors.Wrap(err, "Unable to initialize CRON subsystem (worker)")
		return err
	}
	cronTab.Start()

	_log.Debugf("Cron: initialized")

	// =====================================================================================================================
	//██████╗ ██████╗ ██████╗  ██████╗
	//██╔════╝ ██╔══██╗██╔══██╗██╔════╝
	//██║  ███╗██████╔╝██████╔╝██║
	//██║   ██║██╔══██╗██╔═══╝ ██║
	//╚██████╔╝██║  ██║██║     ╚██████╗
	//╚═════╝ ╚═╝  ╚═╝╚═╝      ╚═════╝

	//--------------------------------------------------------------------------------------
	// goroutine for the standard grpc server
	chanErrGrpc := make(chan error, 1)

	go siloletGrpc.Init(chanErrGrpc)

	//██╗   ██╗███████╗██████╗ ███████╗██╗ ██████╗ ███╗   ██╗
	//██║   ██║██╔════╝██╔══██╗██╔════╝██║██╔═══██╗████╗  ██║
	//██║   ██║█████╗  ██████╔╝███████╗██║██║   ██║██╔██╗ ██║
	//╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██║██║   ██║██║╚██╗██║
	//╚████╔╝ ███████╗██║  ██║███████║██║╚██████╔╝██║ ╚████║
	//╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═══╝

	_log.Debugf("asdlet: Version %s started", Version)

	//--------------------------------------------------------------------------------------
	// vait for channels and return eventual errors
	var errGrpc error

	select {
	case errGrpc = <-chanErrGrpc:
		_log.Error(errGrpc)
		return errGrpc
	}

}

func worker() {
	return
}
