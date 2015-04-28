// router base on reflect
// we would create a new controller for
// every http request
package figo

import (
	"reflect"
)

type Router struct {
	routerMap map[string]reflect.Type
}

func NewRouter() *Router {
	re := &Router{}

	re.routerMap = make(map[string]reflect.Type, 10)

	return re
}

func (this *Router) AddController(path string, c ControllerInterface) {
	this.routerMap[path] = reflect.TypeOf(c)
}

func (this *Router) GetController(path string) ControllerInterface {

	if ct, exists := this.routerMap[path]; exists {
		c := reflect.New(ct)
		return c.Interface().(ControllerInterface)
	}

	return nil
}
