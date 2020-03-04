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
var conn *grpc.ClientConn

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

	s := single.New("asd-client")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		log.Println("another instance of the app is already running, exiting")
		return err
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		log.Println("failed to acquire exclusive app lock: %v", err)
		return err
	}
	defer s.TryUnlock()

	// =========================================================================
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
		_log.Debugf("Debug: listening on: %s", _cfg.Web.DebugHost+":"+_cfg.Web.DebugPort)
		_log.Debugf("Debug: listener closed : %v", http.ListenAndServe(_cfg.Web.DebugHost+":"+_cfg.Web.DebugPort, http.DefaultServeMux))
	}()
	_log.Infof("Debug: initialized")

	// =========================================================================
	//███╗   ███╗███████╗████████╗██████╗ ██╗ ██████╗███████╗
	//████╗ ████║██╔════╝╚══██╔══╝██╔══██╗██║██╔════╝██╔════╝
	//██╔████╔██║█████╗     ██║   ██████╔╝██║██║     ███████╗
	//██║╚██╔╝██║██╔══╝     ██║   ██╔══██╗██║██║     ╚════██║
	//██║ ╚═╝ ██║███████╗   ██║   ██║  ██║██║╚██████╗███████║
	//╚═╝     ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝ ╚═════╝╚══════╝

	// /metrics - Added to the metrics handler

	go func() {
		_log.Debugf("Metrics: listening on %s", _cfg.Web.MetricsHost+":"+_cfg.Web.MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		_log.Debugf("Metrics: listener closed : %v", http.ListenAndServe(_cfg.Web.MetricsHost+":"+_cfg.Web.MetricsPort, nil))
	}()
	_log.Infof("Metrics: initialized")

	//██████╗ ██████╗ ██████╗  ██████╗     ██████╗██╗     ██╗███████╗███╗   ██╗████████╗
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

	// wain intefinitely on the main goroutine
	select {}
}
