package main

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/types"
	_ "database/sql"
	_ "expvar" // Register the expvar handlers
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/gen2brain/beeep"
	"github.com/inconshreveable/go-update"
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
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

	// wain intefinitely on the main goroutine
	select {}
}

//██╗    ██╗ ██████╗ ██████╗ ██╗  ██╗███████╗██████╗
//██║    ██║██╔═══██╗██╔══██╗██║ ██╔╝██╔════╝██╔══██╗
//██║ █╗ ██║██║   ██║██████╔╝█████╔╝ █████╗  ██████╔╝
//██║███╗██║██║   ██║██╔══██╗██╔═██╗ ██╔══╝  ██╔══██╗
//╚███╔███╔╝╚██████╔╝██║  ██║██║  ██╗███████╗██║  ██║
//╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
func worker() {

	// PHASE 1: lock a mutex, these operations are not to be executed concurrently
	_lock.Lock()

	// notify
	_log.Debug("Worker(spawned by internal CRON): starting operations...")

	// PHASE 2: retrieve gp info from postgres db

	dbGps := []types.DbGp{}

	// query for all the gp present
	err := _db.Select(&dbGps, `
	SELECT 
	  * 
	FROM 
	  public.v_utenti
	WHERE
	public.v_utenti.tipo_utente = 'T'
	AND
	public.v_utenti.userid not like '%DEMO%'
	`)

	if err != nil {
		_log.Debugf("DB gp fetch failed: %s", err)
		return
	} else {
		// PHASE 3:

		// prepare Jobs object

		var apiGps api.Gps

		for _, dbGp := range dbGps {
			_log.Debugf("gp:%s", dbGp.Name)
			var apiGp api.Gp
			apiGp.GpId = dbGp.GpId
			apiGps.Gp = append(apiGps.Gp, &apiGp)
		}

		var conn *grpc.ClientConn
		conn, err = grpc.Dial(_cfg.Server.Address+":7777", grpc.WithInsecure())
		if err != nil {
			_log.Errorf("error dialing gRPC server: %s", err)
			_log.Debugf("Retrying again in next service interval")
		} else {
			defer conn.Close()
			c := api.NewServerClient(conn)
			apiJobs, err := c.GetJobs(context.Background(), &apiGps)
			if err != nil {
				_log.Debugf("Error when calling gRPC server: %s", err)
				_log.Debugf("Retrying again in next service interval")
			} else {
				_log.Debugf("Response from server")
				_log.Debugf("Contains %v jobs", len(apiJobs.Job))
				for i, job := range apiJobs.Job {
					// get current time
					start := time.Now()
					_log.Debugf("Job #%v", i+1)
					_log.Debugf("Job_ID=%s, Type=%s, GpId=%s", job.IdJob, job.Type, job.GpId)
					// prepare query

					// determine the right query type

					// query for all the gp present
					var gpMid = findGpMId(dbGps, job.GpId)
					var gpName = findGpName(dbGps, job.GpId)
					// if we have the gp Millewin internal id
					if gpMid != "" {
						content := strings.TrimRight(job.Content, "\r\n")
						_log.Debugf("Querying with gpMid=%s", gpMid)
						rows, err := _db.Queryx(content, gpMid)
						if err != nil {
							_log.Errorf("error in query derived from gRPC call")
							_log.Error(err)
						} else {
							//create or open a log file on the desktop
							// var for *file
							var file *os.File
							// determine home dir
							home, err := homedir.Dir()
							if err != nil {
								_log.Errorf("unable to find home dir")
								_log.Error(err)
							} else {
								// crate full path
								logFilePath := home + "\\Desktop\\asd.log"
								//_log.Debugf("Logfile path=%s", logFilePath)
								file, err = createOrOpen(logFilePath)
								// stat the file
								if err != nil {
									_log.Errorf("unable to create and open existing file on user desktop")
									_log.Error("error:", err)
								} else {
									defer file.Close()
								}
							}
							for rows.Next() {
								results := make(map[string]interface{})
								err = rows.MapScan(results)
								if err != nil {
									_log.Errorf("error in mapscan")
								} else {
									// sort keys, as in pg the same query can give back the same fields
									// in a totally different order, and we need same order

									// create & sort an array with our map keys
									sortedKeys := make([]string, 0, len(results))

									for k := range results {
										sortedKeys = append(sortedKeys, k)
									}
									sort.Strings(sortedKeys)

									var apiResult api.Result
									//for name, value := range results {

									for _, str := range sortedKeys {
										_log.Debugf("name:%s, value:%v", str, results[str])
										var apiEntry api.Entry
										apiEntry.Name = str
										strValue := fmt.Sprintf("%v", results[str])
										apiEntry.Value = strValue
										apiResult.Entries = append(apiResult.Entries, &apiEntry)
									}
									apiResult.GpId = job.GpId
									apiResult.GpName = gpName
									apiResult.Type = job.Type
									apiResult.IdJob = job.IdJob
									elapsed := time.Since(start)
									apiResult.Elapsed = elapsed.Milliseconds()
									// call gRPC method to persist
									apiStatus, err := c.PutResult(context.Background(), &apiResult)
									if err != nil {
										_log.Error(err)
										_log.Errorf("An error occurred in gRPC call PutResult: %s", apiStatus.Message)
									} else {
										Notification("Estrazione trasmessa con successo: " + job.IdJob + "\n" + "Medico: " + gpName)
									}
									// then also write to logfile if we have a correct handler
									if file != nil {
										//write result header
										var str string
										str = "RAPPORTO ESTRAZIONE\n"
										str += "Medico:" + apiResult.GpName + " --- "
										str += "Codice medico:" + apiResult.GpId + " --- "
										str += "Lavoro:" + apiResult.IdJob + " --- "
										str += "Data:" + time.Now().String() + " --- "
										str += "Tempo impiegato (in ms):" + strconv.FormatInt(apiResult.Elapsed, 10) + "\n"
										for _, entry := range apiResult.Entries {
											str += entry.Name + ", "
										}
										str += "\n"
										for _, entry := range apiResult.Entries {
											str += entry.Value + ", "
										}
										str += "\n-------------------------------------------------\n"
										_, err = file.WriteString(str)
										if err != nil {
											_log.Error(err)
											_log.Errorf("error writing to user's desktop logfile")
										}
									}

								}
							}

						}

					}
					// an empty result
					//dbRes := types.DbResult{}
					// fill it with the right data

					// and send it back to the server as a gRPC request
				}
			}
		}
	}

	// PHASE 11: remove lock
	_lock.Unlock()

}

//func connector() {
//	// if not present tries to establish a NATS connection
//	_log.Debug("connector (spawned by cron)")
//
//	// Connect to a server
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		_log.Debugf("Unable to connect to local NATS server?")
//		return
//	}
//
//	// Simple Async Subscriber
//	_, err = nc.Subscribe("exit", func(m *nats.Msg) {
//		fmt.Printf("Received a message: %s\n", string(m.Data))
//		_log.Debugf("Received a message in 'exit' channel")
//		_log.Debugf("Exiting")
//		os.Exit(0)
//	})
//	if err != nil {
//		_log.Debugf("Unable to subscribe to local NATS server on topic 'exit'")
//		return
//	}
//}

//                  _       _
//  _   _ _ __   __| | __ _| |_ ___ _ __
// | | | | '_ \ / _` |/ _` | __/ _ \ '__|
// | |_| | |_) | (_| | (_| | ||  __/ |
//  \__,_| .__/ \__,_|\__,_|\__\___|_|
//       |_|
func updater() {

	_log.Infof("updater() goroutine started")

	_log.Debugf("Trying to update from version: %s", Version)
	// we check a secondary file containing the build number

	baseUrl := "http://asd.avero.it/dist/" + "windows" + "_" + "amd64" + "/"

	// we fetch app own build number
	url := baseUrl + "info.yaml"
	resp, err := http.Get(url)
	if err != nil {
		_log.Errorf("Unable to fetch version file url")
		return
	}

	// create a pointer to an empty UpdateInfo
	updateInfo := types.UpdateInfo{}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_log.Errorf("error: %v", err)
		return
	}

	err = yaml.Unmarshal(bytes, &updateInfo)
	if err != nil {
		_log.Errorf("error: %v", err)
		return
	}

	//_log.Debugf("parsed:%#v", updateInfo)

	_log.Debugf("Fileserver version: %s", updateInfo.Version)

	if updateInfo.Version != Version {

		//	proceed with the update
		_log.Debugf("Updating to version:%s", updateInfo.Version)

		url := baseUrl + "asd-service.exe"
		resp, err := http.Get(url)
		if err != nil {
			_log.Errorf("Unable to fetch update url")
			return
		}
		defer resp.Body.Close()
		opts := update.Options{}
		err = update.Apply(resp.Body, opts)
		if err != nil {
			_log.Error("Unable to update executable file")
			println(err)
			return
		}
		_log.Infof("Executable has been updated. Relaunching...")
		restart()
	} else {
		_log.Debugf("No need to update")
	}
}

//██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗██╗   ██╗██████╗ ██████╗  █████╗ ████████╗███████╗███████╗
//██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝██║   ██║██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝██╔════╝██╔════╝
//██║     ███████║█████╗  ██║     █████╔╝ ██║   ██║██████╔╝██║  ██║███████║   ██║   █████╗  ███████╗
//██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ ██║   ██║██╔═══╝ ██║  ██║██╔══██║   ██║   ██╔══╝  ╚════██║
//╚██████╗██║  ██║███████╗╚██████╗██║  ██╗╚██████╔╝██║     ██████╔╝██║  ██║   ██║   ███████╗███████║
//╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚══════╝
func checkUpdates() {

	_log.Infof("checkUpdates() goroutine started")

	_log.Debugf("Trying to update from version: %s", Version)
	// we check a secondary file containing the build number

	baseUrl := "http://asd.avero.it/dist/" + "windows" + "_" + "amd64" + "/"

	// we fetch app own build number
	url := baseUrl + "info.yaml"
	resp, err := http.Get(url)
	if err != nil {
		_log.Errorf("Unable to fetch version file url")
		return
	}

	// create a pointer to an empty UpdateInfo
	updateInfo := types.UpdateInfo{}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_log.Errorf("error: %v", err)
		return
	}

	err = yaml.Unmarshal(bytes, &updateInfo)
	if err != nil {
		_log.Errorf("error: %v", err)
		return
	}

	_log.Debugf("Fileserver build: %v", updateInfo.Build)

	if updateInfo.Build > Build {
		_log.Debugf("Proceeding to update")

		Notification("E' disponibile una versione piu' aggiornata.")
		//open update url in browser
		open.Start("http://asd.avero.it/download/asd.exe")

	} else {
		_log.Debugf("No need to update")
	}
}

//██╗   ██╗████████╗██╗██╗     ██╗████████╗██╗███████╗███████╗
//██║   ██║╚══██╔══╝██║██║     ██║╚══██╔══╝██║██╔════╝██╔════╝
//██║   ██║   ██║   ██║██║     ██║   ██║   ██║█████╗  ███████╗
//██║   ██║   ██║   ██║██║     ██║   ██║   ██║██╔══╝  ╚════██║
//╚██████╔╝   ██║   ██║███████╗██║   ██║   ██║███████╗███████║
//╚═════╝    ╚═╝   ╚═╝╚══════╝╚═╝   ╚═╝   ╚═╝╚══════╝╚══════╝
func restart() {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	os.StartProcess(os.Args[0], []string{"", "test"}, procAttr)
	os.Exit(0)
}

func findGpMId(dbGp []types.DbGp, gpCode string) string {
	for _, gp := range dbGp {
		if gp.GpId == gpCode {
			return gp.UserID
		}
	}
	return ""
}

func findGpName(dbGp []types.DbGp, gpCode string) string {
	for _, gp := range dbGp {
		if gp.GpId == gpCode {
			return gp.Name + " " + gp.Surname
		}
	}
	return ""
}

func Notification(message string) error {
	err := beeep.Alert("asd", message, "assets/warning.png")
	return err
}

func createOrOpen(path string) (*os.File, error) {
	file, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_log.Errorf("Errore durante l'accesso al file di log: %s", err.Error())
		return nil, err
	} else {
		return file, nil
	}
}
