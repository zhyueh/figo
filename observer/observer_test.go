package observer

import (
	"fmt"
	"testing"
	"time"
)

type ObServerORMData struct {
	Orm int
}

func (this *ObServerORMData) Init() {
	//init orm
	this.Orm = 1
	fmt.Println("in ormdata do")
}

type OrderCancelData struct {
	ObServerORMData
	OrderId    int
	UserId     int
	CancelTime *time.Time
}

type OrderCancelLogHandler struct {
	OrderCancelData //base on OrderCancelData means listening this data
}

func (this *OrderCancelLogHandler) Do() {
	fmt.Println(this.UserId, " cancel an order, orm:", this.Orm)
}

type OrderCancelMessageHandler struct {
	OrderCancelData
}

func (this *OrderCancelMessageHandler) Do() {
	fmt.Println("message handler time", this.CancelTime)
}

// observer data
type OrderCompleteData struct {
	ObServerData
	OrderId int
}

// the first observer handler for OrderCompleteData
type OrderCompleteMessageHandler struct {
	OrderCompleteData
}

func (this *OrderCompleteMessageHandler) Do() {
	fmt.Println("order complete! message handler ", this.OrderId)
}

// the second observer handler for OrderCompleteData
type OrderCompleteLogHandler struct {
	OrderCompleteData
}

func (this *OrderCompleteLogHandler) Do() {
	fmt.Println("order complete! log handler ", this.OrderId)
}

func TestRegisterHandler(t *testing.T) {
	ob := NewObServerCenter()

	// register a handler twice
	ob.RegisterHandler(new(OrderCompleteMessageHandler))
	ob.RegisterHandler(new(OrderCompleteMessageHandler))
	ob.RegisterHandler(new(OrderCompleteLogHandler))
	ob.PrintHandlers()
	if ob.NumHandlerForType(new(OrderCompleteData)) != 2 {
		t.Fatal("register handler fail")
	}
}

func TestHandleData(t *testing.T) {
	ob := NewObServerCenter()
	/*
		tmp := new(OrderCompleteLogHandler)
		tmp.Do()
	*/

	ob.RegisterHandler(new(OrderCompleteMessageHandler))
	ob.RegisterHandler(new(OrderCompleteLogHandler))

	data := new(OrderCompleteData)
	data.OrderId = 4567

	ob.HandleData(data)
}

func TestGlobalCenter(t *testing.T) {

	RegisterHandler(new(OrderCompleteMessageHandler))
	RegisterHandler(new(OrderCompleteLogHandler))
	RegisterHandler(new(OrderCancelLogHandler))
	RegisterHandler(new(OrderCancelMessageHandler))

	data := new(OrderCompleteData)
	data.OrderId = 999
	HandleData(data)

	cdata := new(OrderCancelData)
	cdata.OrderId = 111
	cdata.UserId = 222
	now := time.Now()
	cdata.CancelTime = &now
	HandleData(cdata)

	globalObServerCenter.PrintHandlers()
}
