package main

import (
	"asd/common/helpers"
	"asd/common/types"
	"encoding/json"
	_ "expvar" // Register the expvar handlers
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/gen2brain/beeep"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/marcsauter/single"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	"gopkg.in/stash.v1"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"strconv"
	"strings"
	"sync"
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
		Name     string `conf:"default:hx"`
	}
	Telegram struct {
		Mode string `conf:"default:webhook"` //possible modes "webhook" and "tcp"
		Hook string `conf:"default:https://hx.apps.avero.it"`
		Port string `conf:"default:8000"`
	}
	Log struct {
		Verbose bool `conf:"default:false"`
	}
	FileCache struct {
		Path string `conf:"default:/var/cache/hx"`
	}
}

// =====================================================================================================================
// version and build info

var Build int
var Version string

// =====================================================================================================================
// module wide globals (logging, locks, db, etc)
var _log *logging.Logger
var _verbose bool

var _db *sqlx.DB

var dbItems []types.DbItem
var dbItemsLock sync.RWMutex

var sbItems types.SbItems
var sbItemsLock sync.RWMutex

var bot *tgbotapi.BotAPI

var new chan types.ServerInfo

var fileCache *stash.Cache

// =====================================================================================================================
// metrics globals
var (
	_workerOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hx_worker_ops_total",
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

	if err := conf.Parse(os.Args[1:], "HX", &_cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("HX", &_cfg)
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

	s := single.New("hx")
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
	//      _ _
	//   __| | |__
	//  / _` | '_ \
	// | (_| | |_) |
	//  \__,_|_.__/
	//

	var connectString string
	if _cfg.Db.Instance != "" {
		connectString = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", _cfg.Db.User, _cfg.Db.Password, _cfg.Db.Instance, _cfg.Db.Name)
	} else {
		connectString = _cfg.Db.User + ":" + _cfg.Db.Password + "@(" + _cfg.Db.Host + ":" + _cfg.Db.Port + ")/" + _cfg.Db.Name + "?parseTime=true"
	}
	//_log.Infof(connectString)
	var err error
	_db, err = sqlx.Connect("mysql", connectString)
	if err != nil {
		return errors.Wrap(err, "connecting to backend db")
	}
	_log.Infof("Db: initialized")

	// =====================================================================================================================
	//
	//      _                            _
	//  ___| |__   __ _ _ __  _ __   ___| |___
	// / __| '_ \ / _` | '_ \| '_ \ / _ \ / __|
	//| (__| | | | (_| | | | | | | |  __/ \__ \
	// \___|_| |_|\__,_|_| |_|_| |_|\___|_|___/

	new = make(chan types.ServerInfo, 100)
	_log.Infof("Channels: initialized")

	// =====================================================================================================================
	//  _           _
	// | |__   ___ | |_
	// | '_ \ / _ \| __|
	// | |_) | (_) | |_
	// |_.__/ \___/ \__|

	bot, err = tgbotapi.NewBotAPI("1021934418:AAFwdPiIkJS1iJESJJGSRnk5jIXA0jP_pa0")
	if err != nil {
		return errors.Wrap(err, "listening on Telegram tcp")
	}

	if _cfg.Telegram.Mode == "tcp" {
	}

	if _cfg.Telegram.Mode == "webhook" {
		bot.Debug = true

		_log.Infof("Authorized on account %s", bot.Self.UserName)

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(_cfg.Telegram.Hook + "/" + bot.Token))
		if err != nil {
			return errors.Wrap(err, "setting Telegram webhook")
		}

		go http.ListenAndServe(":"+_cfg.Telegram.Port, nil)
	}

	_log.Infof("Bot: initialized")
	_log.Debugf("Bot: authorized on account %s", bot.Self.UserName)

	// =====================================================================================================================
	// STASH

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
	//                     _
	// __   _____ _ __ ___(_) ___  _ __
	// \ \ / / _ \ '__/ __| |/ _ \| '_ \
	//  \ V /  __/ |  \__ \ | (_) | | | |
	//   \_/ \___|_|  |___/_|\___/|_| |_|

	_log.Infof("HX: Version %s started", Version)

	// =====================================================================================================================
	// execute on startup
	worker()

	go replayer(bot)

	go poster(bot)

	// =====================================================================================================================
	// wain intefinitely on the main goroutine
	select {}
}

func poster(bot *tgbotapi.BotAPI) {
	for {
		serverInfo := <-new
		msg := tgbotapi.NewMessage(serverInfo.DbItem.ChatId, "")
		msg.ParseMode = "markdown"

		if serverInfo.Event == "VARIATION" {
			msg.Text = "PRICE VARIATION\n"
			msg.Text += serverInfoDetails(serverInfo.Server, serverInfo.DbItem, true)
		}

		if serverInfo.Event == "VANISHED" {
			msg.Text = "SERVER VANISHED\n"
			msg.Text += "ID=" + fmt.Sprint(serverInfo.DbItem.ServerId) + "\n"
		}
		bot.Send(msg)
	}
}

func replayer(bot *tgbotapi.BotAPI) {

	var updates <-chan tgbotapi.Update
	var err error

	if _cfg.Telegram.Mode == "tcp" {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			errors.Wrap(err, "bot: unable to get updates")
		}
	}

	if _cfg.Telegram.Mode == "webhook" {
		updates = bot.ListenForWebhook("/" + bot.Token)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			_log.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			words := strings.Fields(update.Message.Text)
			switch update.Message.Command() {
			case "help":
				msg.Text = "HELP\n/watch <id>\n/forget <id>\n/list\n/astiastiieri"
			case "watch":
				if len(words) == 2 {
					//correct number of elements
					arg := words[1]
					msg.Text = watch(arg, update.Message.Chat.ID, update.Message.From.UserName)
				} else {
					msg.Text = "wrong number of parameters"
				}
			case "forget":
				if len(words) == 2 {
					//correct number of elements
					arg := words[1]
					msg.Text = forget(arg)
				} else {
					msg.Text = "wrong number of parameters"
				}
			case "list":
				if len(words) == 1 {
					//correct number of elements
					msg.Text = list()
				} else {
					msg.Text = "wrong number of parameters"
				}
			case "astiastiieri":
				msg.Text = "Puppa"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = "markdown"
			bot.Send(msg)
		}
	}

}

func list() string {
	msg := ""
	err := selectItems()
	if err != nil {
		_log.Error("Error selecting db items (list)")
	} else {
		sbItemsLock.RLock()
		dbItemsLock.RLock()
		msg += "Elenco server\n"
		for _, sbitem := range sbItems.Server {
			for _, dbitem := range dbItems {
				if dbitem.ServerId == sbitem.Key {
					msg += serverInfoDetails(sbitem, dbitem, false)
					//msg += sbitem.Description + "/n"
				}
			}
		}
		sbItemsLock.RUnlock()
		dbItemsLock.RUnlock()
	}
	return msg
}

func serverInfoDetails(sbItem types.Server, dbItem types.DbItem, diff bool) string {
	msg := "\n"
	sKey := fmt.Sprint(sbItem.Key)
	sCPUBenchmark := fmt.Sprint(sbItem.CPUBenchmark)
	fPrice, err := strconv.ParseFloat(sbItem.Price, 64)
	if err != nil {
		fPrice = 0.0
		_log.Errorf("Error parsing price from sbItem struct")
	}
	fPriceVat := fPrice * 1.22
	sPriceVat := fmt.Sprintf("%.2f", fPriceVat)
	msg += "*WATCHER:* " + dbItem.UserName + "\n"
	msg += "*ID:* " + sKey + ", *NAME:* " + sbItem.Name + ", *CPU:* " + sbItem.CPU + ", *CPUMARK:* " + sCPUBenchmark + "\n"
	msg += "*PRICE:* " + sbItem.Price + ", *PRICE*(inc. VAT): " + sPriceVat + "\n"

	if diff {
		msg += "Previous price:\n"
		oldPrice := dbItem.Price
		fPrice, err := strconv.ParseFloat(oldPrice, 64)
		if err != nil {
			fPrice = 0.0
			_log.Errorf("Error parsing price from sbItem struct")
		}
		fPriceVat := fPrice * 1.22
		sPriceVat := fmt.Sprintf("%.2f", fPriceVat)
		msg += "*PRICE:* " + oldPrice + ", *PRICE*(inc. VAT): " + sPriceVat + "\n"
	}

	msg += "[ORDER](http://scnet.link/hetzner.html?id=" + sKey + ")\n"
	return msg
}

func watch(arg string, chatId int64, userName string) string {
	dbItemsLock.Lock()
	serverId64, err := strconv.ParseInt(arg, 10, 64)
	serverId := int(serverId64)
	price, err := fetchPrice(serverId)
	_, err = _db.Queryx("REPLACE INTO items VALUES (?,?,?,?,?);", 0, serverId, price, chatId, userName)
	if err != nil {
		msg := "error in REPLACE query (watch)"
		_log.Error(msg, err)
		dbItemsLock.Unlock()
		return msg
	} else {
		msg := "new server added to watch list"
		_log.Info(msg, err)
		dbItemsLock.Unlock()
		return msg
	}
}

func forget(arg string) string {
	dbItemsLock.Lock()
	serverId64, err := strconv.ParseInt(arg, 10, 64)
	serverId := int(serverId64)
	//price, err := fetchPrice(serverId)
	_, err = _db.Queryx("DELETE FROM items WHERE server_id = ?;", serverId)
	if err != nil {
		msg := "error in DELETE query (forget)"
		_log.Error(msg, err)
		dbItemsLock.Unlock()
		return msg
	} else {
		msg := "server removed from watch list"
		_log.Info(msg, err)
		dbItemsLock.Unlock()
		return msg
	}
}

func fetchPrice(serverId int) (float64, error) {
	sbItemsLock.RLock()
	var price float64
	var found bool
	var err error
	for _, item := range sbItems.Server {
		if item.Key == serverId {
			price, err = strconv.ParseFloat(item.Price, 64)
			if err != nil {
				_log.Error("error parsing string (price) to float64")
				continue
			} else {
				found = true
				continue
			}
		}
	}
	sbItemsLock.RUnlock()

	if found {
		return price, nil
	} else {
		return 0.0, err
	}
}

func selectItems() error {
	// global object, one writer a time, or n readers together
	dbItemsLock.Lock()
	dbItems = dbItems[:0]
	// query for all the items present
	err := _db.Select(&dbItems, `
	SELECT * FROM items`)
	if err != nil {
		_log.Debugf("DB item fetch failed: %s", err)
		dbItemsLock.Unlock()
		return err
	}

	// done wiriting
	dbItemsLock.Unlock()

	return nil
}

func worker() {

	// increment metrics counter
	_workerOps.Inc()

	// get time
	start := time.Now()

	//_log.Debug(bot)

	// notify
	_log.Debug("Worker: resuming operations...")

	err := selectItems()
	if err != nil {
		_log.Error("Error selecting db items")
	} else {
		// PHASE 3: fetch json from api, marshal, into sbItems (hetzner server bidding items)

		var body []byte

		// json data
		url := "https://www.hetzner.com/a_hz_serverboerse/live_data.json"
		res, err := http.Get(url)
		if err != nil {
			_log.Errorf("failed to open JSON resource with the following error")
			_log.Error(err)
		} else {
			body, err = ioutil.ReadAll(res.Body)
			if err != nil {
				_log.Errorf("failed to download JSON resource with the following error")
				_log.Error(err)
				return
			}
			// update file cache with []byte
			err = fileCache.Put("sbItems", body)
			if err != nil {
				_log.Errorf("failed to store body in filecache")
				_log.Error(err)
				return
			}
		}

		if len(body) == 0 {
			// it network fetch fails, use the cache
			_log.Notice("trying to acquire old version from filecache")
			reader, err := fileCache.Get("sbItems")
			if err != nil {
				_log.Errorf("failed to open file in filecache")
				_log.Error(err)
				return
			}

			body, err = ioutil.ReadAll(reader) // => []byte("Hello, world!\n")
			if err != nil {
				_log.Errorf("failed to read body from filecache")
				_log.Error(err)
				return
			} else {
				_log.Noticef("sbItems read from filecache")
				reader.Close()
			}
		}

		// global object, one writer a time, or n readers together
		sbItemsLock.Lock()

		// unmarshael to sbItems
		err = json.Unmarshal(body, &sbItems)

		// done wiriting
		sbItemsLock.Unlock()

		if err != nil {
			_log.Errorf("failed to unmarshal JSON resource with the following error")
			_log.Error(err)
			return
		}

		sbItemsLock.RLock()
		dbItemsLock.RLock()
		var found bool

		for _, dbItem := range dbItems {
			found = false
			for _, sbItem := range sbItems.Server {
				if dbItem.ServerId == sbItem.Key {
					found = true
					//if same id compare price
					if dbItem.Price != sbItem.Price {
						// if price has changed
						// log it
						_log.Infof("Price change detected")
						_log.Infof("Key: %v desc: %s", sbItem.Key, sbItem.Name)
						_log.Infof("new price: %s", sbItem.Price)
						_log.Infof("old price: %s", dbItem.Price)
						var oldPrice, newPrice float64
						oldPrice, err = strconv.ParseFloat(dbItem.Price, 64)
						newPrice, err = strconv.ParseFloat(sbItem.Price, 64)
						_log.Infof("...change: %v", oldPrice-newPrice)
						// then update the db
						_log.Debugf("persisting record...")
						_, err := _db.Queryx("UPDATE items SET price = ? WHERE server_id = ?;", sbItem.Price, sbItem.Key)
						if err != nil {
							_log.Error("error in update query:", err)
						} else {
							_log.Infof("price persisted")
							serverInfo := types.ServerInfo{"VARIATION", sbItem, dbItem}
							new <- serverInfo
						}
					}
				}
			}
			if !found {
				serverInfo := types.ServerInfo{"VANISH", types.Server{}, dbItem}
				new <- serverInfo
			}
		}

		dbItemsLock.RUnlock()
		sbItemsLock.RUnlock()

		/*		var keys []int
				for _, item := range sbItems.Server {
					keys = append(keys, item.Key)
				}

				sort.Ints(keys)
				var last int
				for _, key := range keys {
					if last == key {
						_log.Info("duplicate key:", key)
					}
					last = key
					//_log.Info("key:", key)
				}
		*/
		// if they match and same price, no action

		// if they match and new price -> telegram report new price

		// if no match -> telegram report missing server

	}

	// show elapsed time
	elapsed := time.Since(start)
	_log.Debugf("Worker took %s", elapsed)
}

func Notification(message string) error {
	err := beeep.Alert("HX", message, "assets/warning.png")
	return err
}
