package main

import (
	"asd/common/api"
	"asd/common/helpers"
	_ "expvar"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"gopkg.in/stash.v1"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

// =====================================================================================================================
// configuration structure
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
	Db struct {
		User     string `conf:"default:joegalaxy"`
		Password string `conf:"default:1hugesoul,noprint"`
		Host     string `conf:"default:127.0.0.1"`
		Instance string
		Port     string `conf:"default:3306"`
		Name     string `conf:"default:asd"`
	}
	Telegram struct {
		Mode string `conf:"default:webhook"` //possible modes "webhook" and "tcp"
		Hook string `conf:"default:https://asd.apps.avero.it"`
		Port string `conf:"default:8000"`
	}
	Log struct {
		Verbose bool `conf:"default:false"`
	}
	FileCache struct {
		Path string `conf:"default:/var/cache/asd"`
	}
}

// =====================================================================================================================
// version and build info

var Build int
var Version string

// =====================================================================================================================
// Server represents the gRPC server

type Server struct {
}

// =====================================================================================================================
// module wide globals (logging, locks, db, etc)
var _log *logging.Logger
var _verbose bool

var fileCache *stash.Cache

// =====================================================================================================================
// metrics globals
var (
	_workerOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "asd_worker_ops_total",
		Help: "The total number worker was spawned by internal cron subsys",
	})
)

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
	Build = 0_3_0

	// =====================================================================================================================
	//                   __ _
	//	 ___ ___  _ __  / _(_) __ _
	// 	/ __/ _ \| '_ \| |_| |/ _` |
	// | (_| (_) | | | |  _| | (_| |
	// 	\___\___/|_| |_|_| |_|\__, |
	//                         |___/

	if err := conf.Parse(os.Args[1:], "ASD", &_cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("ASD", &_cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =====================================================================================================================
	//	_
	// | | ___   __ _ ___
	// | |/ _ \ / _` / __|
	// | | (_) | (_| \__ \
	// |_|\___/ \__, |___/
	//          |___/

	_log = helpers.InitLogs(_cfg.Log.Verbose)
	_log.Info("Log: initialized")

	// =====================================================================================================================
	//    	_                   _
	//  ___(_) __ _ _ __   __ _| |___
	// / __| |/ _` | '_ \ / _` | / __|
	// \__ \ | (_| | | | | (_| | \__ \
	// |___/_|\__, |_| |_|\__,_|_|___/
	//        |___/
	helpers.InitSignals()
	_log.Infof("Signals: initialized")

	// =====================================================================================================================
	//              _
	//  _   _ _ __ (_) __ _ _   _  ___
	// | | | | '_ \| |/ _` | | | |/ _ \
	// | |_| | | | | | (_| | |_| |  __/
	//  \__,_|_| |_|_|\__, |\__,_|\___|
	//                   |_|

	s := single.New("asd")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		_log.Errorf("another instance of the app is already running, exiting")
		return err
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		_log.Errorf("failed to acquire exclusive app lock: %v", err)
		return err
	}
	defer s.TryUnlock()
	_log.Infof("Unique: initialized")

	// =========================================================================
	//      _      _
	//   __| | ___| |__  _   _  __ _
	//  / _` |/ _ \ '_ \| | | |/ _` |
	// | (_| |  __/ |_) | |_| | (_| |
	//  \__,_|\___|_.__/ \__,_|\__, |
	//                         |___/

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
	//                 _
	//  _ __ ___   ___| |_ _ __(_) ___ ___
	// | '_ ` _ \ / _ \ __| '__| |/ __/ __|
	// | | | | | |  __/ |_| |  | | (__\__ \
	// |_| |_| |_|\___|\__|_|  |_|\___|___/

	// /metrics - Added to the metrics handler

	go func() {
		_log.Debugf("Metrics: listening on %s", _cfg.Web.MetricsHost+":"+_cfg.Web.MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		_log.Debugf("Metrics: listener closed : %v", http.ListenAndServe(_cfg.Web.MetricsHost+":"+_cfg.Web.MetricsPort, nil))
	}()
	_log.Infof("Metrics: initialized")

	// =====================================================================================================================
	//
	//      _                            _
	//  ___| |__   __ _ _ __  _ __   ___| |___
	// / __| '_ \ / _` | '_ \| '_ \ / _ \ / __|
	//| (__| | | | (_| | | | | | | |  __/ \__ \
	// \___|_| |_|\__,_|_| |_|_| |_|\___|_|___/

	_log.Infof("Channels: initialized")

	// =====================================================================================================================
	// STASH
	var err error
	fileCache, err = stash.New(_cfg.FileCache.Path, 10000000, 1000)
	if err != nil {
		errors.Wrap(err, "Unable to initialize file cache, path="+_cfg.FileCache.Path)
		return err
	}

	// =====================================================================================================================
	//   ___ _ __ ___  _ __
	//  / __| '__/ _ \| '_ \
	// | (__| | | (_) | | | |
	//  \___|_|  \___/|_| |_|

	cronTab := cron.New()
	err = cronTab.AddFunc("0 * * * * *", worker) //every minute
	if err != nil {
		errors.Wrap(err, "Unable to initialize CRON subsystem (worker)")
		return err
	}
	cronTab.Start()

	_log.Infof("Cron: initialized")

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
	go func(errGrpc chan<- error) {
		// create a listener on TCP port 7777
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
		if err != nil {
			err = errors.Wrap(err, "trying to listen on tcp")
			errGrpc <- err
			return
		}
		// create a server instance
		s := Server{}
		// create a gRPC server object
		grpcServer := grpc.NewServer()
		// attach the Ping service to the server
		api.RegisterAsdServer(grpcServer, &s)
		// start the server
		_log.Infof("listening for grpc connections on port: 7777")
		if err := grpcServer.Serve(lis); err != nil {
			err = errors.Wrap(err, "failed to serve gRPC")
			errGrpc <- err
			return
		}
	}(chanErrGrpc)

	//--------------------------------------------------------------------------------------
	// goroutine for the grpc<->http reverse proxy
	chanErrHttp := make(chan error, 1)
	go func(errHttp chan<- error) {

		// http server
		port := 8080

		router := mux.NewRouter()

		// pages
		//router.HandleFunc("/readall", readAll).Methods("GET")
		//router.HandleFunc("/create", create).Methods("POST")
		//router.HandleFunc("/update", update).Methods("POST")
		//router.HandleFunc("/delete", delete).Methods("POST")
		//router.HandleFunc("/getallresults", getAllResults).Methods("GET")
		//router.HandleFunc("/getalljobs", getAllJobs).Methods("GET")
		//router.HandleFunc("/getallgps", getAllGps).Methods("GET")

		http.Handle("/", router)
		// start the web server (blocking)
		_log.Infof("listening for http connections on port: %v", port)
		if err := http.ListenAndServe(fmt.Sprint(":", port), router); err != nil {
			err = errors.Wrap(err, "failed to serve http<-->gRPC")
			errHttp <- err
			return
		}
	}(chanErrHttp)

	//--------------------------------------------------------------------------------------
	// vait for channels and return eventual errors
	var errGrpc, errHttp error

	select {
	case errGrpc = <-chanErrGrpc:
		_log.Error(errGrpc)
		return errGrpc
	case errHttp = <-chanErrHttp:
		_log.Error(errHttp)
		return errHttp
	}

	// =====================================================================================================================
	//                     _
	// __   _____ _ __ ___(_) ___  _ __
	// \ \ / / _ \ '__/ __| |/ _ \| '_ \
	//  \ V /  __/ |  \__ \ | (_) | | | |
	//   \_/ \___|_|  |___/_|\___/|_| |_|

	_log.Infof("asd: Version %s started", Version)

	// =====================================================================================================================
	// wain intefinitely on the main goroutine
	select {}
}

func worker() {
	return
}
