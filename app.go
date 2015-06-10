package figo

import (
	"fmt"
	"github.com/zhyueh/figo/cache"
	"github.com/zhyueh/figo/log"
	"github.com/zhyueh/figo/orm"
	"github.com/zhyueh/figo/toolkit"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

type App struct {
	name   string
	config *AppConfig

	routerFunc func(string) ControllerInterface

	accessLogger      *log.DataLogger
	appLogger         *log.DataLogger
	controllerLoggers map[string]*log.DataLogger
	cache             *cache.Cache
	orm               *orm.Orm
}

func NewApp(name string) *App {
	re := &App{}
	re.name = name
	re.config = nil
	re.accessLogger = nil
	re.appLogger = nil
	re.controllerLoggers = make(map[string]*log.DataLogger, 16)
	return re
}

func (this *App) Run(config *AppConfig, routerFunc func(string) ControllerInterface) {
	this.config = config

	if logger, err := log.NewDataLogger(this.config.LogPath, "access", log.TruncateHour, log.LevelInfo, 1000); err != nil {
		panic(err)
	} else {
		this.accessLogger = logger
	}

	if logger, err := log.NewDataLogger(this.config.LogPath, "error", log.TruncateImmediately, log.LevelError, 1000); err != nil {
		panic(err)
	} else {
		this.appLogger = logger
	}

	if cache, err := cache.NewCache(config.CacheAddress, config.CachePassword, config.CacheDB); err != nil {
		this.appLogger.Fatal("cache init error %v", err)
	} else {
		this.cache = cache
	}

	orm, ormerr := orm.NewOrm(
		config.DbType,
		config.DbHost,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbPort,
	)

	if ormerr != nil {
		this.appLogger.Fatal("orm init error %v", ormerr)
	}
	this.orm = orm

	runtime.GOMAXPROCS(runtime.NumCPU())

	toolkit.SignalWatchRegister(this.quit, os.Kill, os.Interrupt, syscall.SIGTERM)
	toolkit.SignalWatchRun()
	//init log for app

	listen, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Port))

	this.routerFunc = routerFunc
	err := http.Serve(listen, this)

	if err != nil {
		this.appLogger.Fatal("%v", err)
	}
}

func (this *App) getControllerLogger(name string) *log.DataLogger {
	if name == "" {
		name = "default"
	}

	if logger, exists := this.controllerLoggers[name]; exists {
		return logger
	}
	logger := this.createControllerLogger(name)
	this.controllerLoggers[name] = logger
	return logger
}

func (this *App) createControllerLogger(name string) *log.DataLogger {
	logger, err := log.NewDataLogger(this.config.LogPath, name, log.TruncateHour, log.LevelInfo, 1000)
	if err != nil {
		this.appLogger.Fatal("%v", err)
	}
	return logger
}

func (this *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	httpStatus := this.safeRun(w, r)

	end := time.Now()
	timeCost := float32(end.UnixNano()-start.UnixNano()) / 1000000.0

	accessMsg := fmt.Sprintf(" %d %s %s (%s) %.2fms", httpStatus, r.Method, r.URL.Path, toolkit.GetRemoteIP(r), timeCost)

	this.accessLogger.Info(accessMsg)

}

func (this *App) safeRun(w http.ResponseWriter, r *http.Request) (httpStatus int) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintln(r.URL.Path, err, ";trace", string(debug.Stack()))
			this.appLogger.Error(msg)
			httpStatus = http.StatusInternalServerError
		}
	}()

	controller := this.routerFunc(r.URL.Path)

	if controller != nil {
		controller.Init(w, r)
		controller.SetLogger(this.getControllerLogger(controller.GetControllerName()))
		controller.SetCache(this.cache)
		controller.SetOrm(this.orm)
		preloadErr := controller.Preload()
		if preloadErr == nil {
			switch controller.GetConnectMode() {
			case HttpMode:
				method := r.Method
				switch method {
				case "POST":
					controller.Post()
				default:
					controller.Get()
				}
				//campatible websocket
			case WebsocketMode:
				webHandler := func(conn *websocket.Conn) {
					controller.SetWebsocketConnection(conn)
					controller.Get()
				}
				s := websocket.Server{Handler: webHandler}
				s.ServeHTTP(w, r)
			}
		}

		controller.Flush()
		httpStatus = http.StatusOK
	} else {
		httpStatus = http.StatusNotFound
	}

	return

}

func (this *App) quit() {
	log.CloseAllLogs()
	this.cache.Close()
	os.Exit(0)
}
