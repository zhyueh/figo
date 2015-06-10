package figo

import (
	"github.com/zhyueh/figo/cache"
	"github.com/zhyueh/figo/log"
	"github.com/zhyueh/figo/orm"
	"golang.org/x/net/websocket"
	"net/http"
)

const (
	HttpMode      = 0
	WebsocketMode = 1
)

type ControllerInterface interface {
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

}

func (this *Controller) Post() {

}
