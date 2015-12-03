package observer

import (
	"fmt"
	"reflect"
)

// all data handler should implement do func
type ObServerDataInterface interface {
	Do()
}

type ObServerData struct {
	//leaving blank
}

func (this *ObServerData) Do() {
	//leaving blank
}

type ObServerCenter struct {
	// store all the handler type
	registerHandler map[reflect.Type][]reflect.Type
}

func (this *ObServerCenter) Init() {
	this.registerHandler = make(map[reflect.Type][]reflect.Type, 10)
}

// checkout the server data type which i is implemented from
func (this *ObServerCenter) getServerDataTypeFromServerDataInterface(i ObServerDataInterface) reflect.Type {
	iType := reflect.ValueOf(i).Elem().Type()
	if iType.NumField() > 0 {
		firstFieldType := iType.Field(0).Type
		if firstFieldType.Kind() == reflect.Struct {
			return firstFieldType
		}
	}

	return iType
}

func (this *ObServerCenter) RegisterHandler(i ObServerDataInterface) {
	handlerType := reflect.ValueOf(i).Elem().Type()
	dataType := this.getServerDataTypeFromServerDataInterface(i)
	var v []reflect.Type
	var ok bool
	if v, ok = this.registerHandler[dataType]; !ok {
		v = make([]reflect.Type, 0)
	}
	for _, ht := range v {
		if ht == handlerType {
			return
		}
	}

	v = append(v, handlerType)
	this.registerHandler[dataType] = v
}

func (this *ObServerCenter) PrintHandlers() {
	fmt.Println(this.registerHandler)
}

func (this *ObServerCenter) handlersForType(i ObServerDataInterface) []reflect.Type {
	dataType := reflect.ValueOf(i).Elem().Type()
	if v, ok := this.registerHandler[dataType]; ok {
		return v
	}
	return []reflect.Type{}
}

func (this *ObServerCenter) NumHandlerForType(i ObServerDataInterface) int {
	return len(this.handlersForType(i))
}

// find out all the register handler for this server data
// and then run the do function after creating handler instances
// todo using chan and go . so that handlers can run in async
func (this *ObServerCenter) HandleData(i ObServerDataInterface) {

	iv := reflect.ValueOf(i).Elem()

	for _, handlerType := range this.handlersForType(i) {
		handlerPtr := reflect.New(handlerType)
		handler := handlerPtr.Elem()

		//copy value from i into handler
		if handler.NumField() > 0 {

			// even it has check the type when register handler
			// but we have to ensure the data's type is the same as handler's
			handlerBaseType := handler.Field(0).Type()
			dataType := iv.Type()
			if handlerBaseType.PkgPath() != dataType.PkgPath() ||
				handlerBaseType.String() != dataType.String() {
				continue
			}

			handler.Field(0).Set(iv)

			//execute do of handler
			doFunc := handlerPtr.MethodByName("Do")
			if doFunc.IsValid() {
				doFunc.Call(nil)
			} else {
				fmt.Println(handlerPtr, "has no do func")
			}
		}

	}
}

func NewObServerCenter() *ObServerCenter {
	re := new(ObServerCenter)
	re.Init()
	return re
}

var globalObServerCenter = NewObServerCenter()

func RegisterHandler(i ObServerDataInterface) {
	globalObServerCenter.RegisterHandler(i)
}

func HandleData(i ObServerDataInterface) {
	globalObServerCenter.HandleData(i)
}
