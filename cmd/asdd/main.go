package main

import (
	"asd/common/api"
	"asd/common/helpers"
	_ "expvar"
	"fmt"
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
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

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

	//██╗      ██████╗  ██████╗ ███████╗
	//██║     ██╔═══██╗██╔════╝ ██╔════╝
	//██║     ██║   ██║██║  ███╗███████╗
	//██║     ██║   ██║██║   ██║╚════██║
	//███████╗╚██████╔╝╚██████╔╝███████║
	//╚══════╝ ╚═════╝  ╚═════╝ ╚══════╝

	_log = helpers.InitLogs(true)
	_log.Debug("Log: initialized")

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

	s := single.New("asdd")
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
	err := cronTab.AddFunc("0 * * * * *", worker) //every minute
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
		_log.Debugf("listening for grpc connections on port: 7777")
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
		_log.Debugf("listening for http connections on port: %v", port)
		if err := http.ListenAndServe(fmt.Sprint(":", port), router); err != nil {
			err = errors.Wrap(err, "failed to serve http<-->gRPC")
			errHttp <- err
			return
		}
	}(chanErrHttp)

	//██╗   ██╗███████╗██████╗ ███████╗██╗ ██████╗ ███╗   ██╗
	//██║   ██║██╔════╝██╔══██╗██╔════╝██║██╔═══██╗████╗  ██║
	//██║   ██║█████╗  ██████╔╝███████╗██║██║   ██║██╔██╗ ██║
	//╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██║██║   ██║██║╚██╗██║
	//╚████╔╝ ███████╗██║  ██║███████║██║╚██████╔╝██║ ╚████║
	//╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═══╝

	_log.Debugf("asd: Version %s started", Version)

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

}

func worker() {
	return
}
