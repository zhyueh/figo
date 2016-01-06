package figo

import (
	"fmt"
	"github.com/zhyueh/figo/cache"
	"github.com/zhyueh/figo/log"
	"github.com/zhyueh/figo/orm"
	"github.com/zhyueh/figo/toolkit"
	"golang.org/x/net/websocket"
	"net/http"
	"reflect"
	"strings"
)

const (
	HttpMode      = 0
	WebsocketMode = 1
)

const (
	SingleRouterMode = 0
	AutoRouterMode   = 1
)

type ControllerInterface interface {
	//phase II
	Init(w http.ResponseWriter, r *http.Request)
	Preload() error
	GetConnectMode() int8
	Get()
	Post()
	Flush()
	SetLogger(*log.DataLogger)
	SetCache(*cache.Cache)
	SetOrm(*orm.Orm)
	SetWebsocketConnection(*websocket.Conn)

	GetControllerName() string

	//phase II
	GetRouteMode() int8
	//Route()
	HandleRequestPathNotFound()
	HandleRequestParseError(error)
}

type Controller struct {
	Req            *Requester
	Resp           *Response
	ControllerName string
	Logger         *log.DataLogger
	Cache          *cache.Cache
	Orm            *orm.Orm
	WebsocketConn  *websocket.Conn
}

func NewController() *Controller {
	re := &Controller{}

	re.ControllerName = ""

	return re
}

func (this *Controller) Preload() error {
	return nil
}

func (this *Controller) GetConnectMode() int8 {
	return HttpMode
}

func (this *Controller) SetWebsocketConnection(conn *websocket.Conn) {
	this.WebsocketConn = conn
}

func (this *Controller) SetCache(cache *cache.Cache) {
	this.Cache = cache
}
func (this *Controller) SetLogger(logger *log.DataLogger) {
	this.Logger = logger
}

func (this *Controller) SetOrm(orm *orm.Orm) {
	this.Orm = orm
}

func (this *Controller) GetControllerName() string {
	return this.ControllerName
}

func (this *Controller) Init(w http.ResponseWriter, r *http.Request) {
	this.Req = NewRequester(r)
	this.Resp = NewResponse(w)

}

func (this *Controller) Flush() {
	this.Resp.Flush()
}

func (this *Controller) Get() {
	fmt.Println("not implement get")
}

func (this *Controller) Post() {

}

func (this *Controller) GetRouteMode() int8 {
	//default single router mode
	return SingleRouterMode
}

func ControllerHandleFunc(this ControllerInterface, httpMethod, path string) {
	if this.GetRouteMode() == SingleRouterMode {
		switch httpMethod {
		case "POST":
			this.Post()
		default:
			this.Get()
		}
	} else {
		autoRouteFunc(this, httpMethod, path)
	}
}

func autoRouteFunc(this ControllerInterface, httpMethod, path string) {
	comps := toolkit.SplitString(path, "/")
	if len(comps) < 2 {
		this.HandleRequestPathNotFound()
		return
	}
	var methodName string
	for _, v := range comps[1:] {
		subComps := strings.Split(v, "-")
		for _, vv := range subComps {
			methodName += strings.Title(vv)
		}
	}
	if len(methodName) == 0 {
		this.HandleRequestPathNotFound()
		return
	}
	methodName = strings.Title(strings.ToLower(httpMethod)) + methodName
	value := reflect.ValueOf(this)
	method := value.MethodByName(methodName)
	if !method.IsValid() {
		this.HandleRequestPathNotFound()
		return
	}
	method.Call(nil)
}

func (this *Controller) HandleRequestParseError(err error) {
	this.Resp.WriteString(fmt.Sprintf("parse err %v", err))
}

func (this *Controller) HandleRequestPathNotFound() {
	this.Resp.WriteString(fmt.Sprintf("path not found %s", this.Req.BaseHttpRequest.URL.Path))
}
