package observer

import (
	"fmt"
	"testing"
)

type OrderCancelData struct {
	ObServerData
	OrderId int
	UserId  int
}

type OrderCancelLogHandler struct {
	OrderCancelData //base on OrderCancelData means listening this data
}

func (this *OrderCancelLogHandler) Do() {
	fmt.Println(this.UserId, " cancel an order")
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
	ob.RegisterHandler(new(OrderCancelLogHandler))

	data := new(OrderCompleteData)
	data.OrderId = 456

	ob.HandleData(data)

	cdata := new(OrderCancelData)
	cdata.OrderId = 777
	cdata.UserId = 888

	ob.HandleData(cdata)

}

func TestGlobalCenter(t *testing.T) {
	RegisterHandler(new(OrderCompleteMessageHandler))
	RegisterHandler(new(OrderCompleteLogHandler))
	RegisterHandler(new(OrderCancelLogHandler))

	data := new(OrderCompleteData)
	data.OrderId = 999

	HandleData(data)

	cdata := new(OrderCancelData)
	cdata.OrderId = 111
	cdata.UserId = 222

	HandleData(cdata)
}
