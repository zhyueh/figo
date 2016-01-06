// router base on reflect
// we would create a new controller for
// every http request
package figo

import (
	//"fmt"
	"github.com/zhyueh/figo/toolkit"
	"reflect"
)

type Router struct {
	routerMap     map[string]reflect.Type
	autoRouterMap map[string]reflect.Type
	routerFunc    func(string) ControllerInterface
}

func NewRouter() *Router {
	re := &Router{}

	re.routerFunc = nil

	re.routerMap = make(map[string]reflect.Type, 10)
	re.autoRouterMap = make(map[string]reflect.Type, 10)

	return re
}

func (this *Router) SetRouterFunc(routerFunc func(string) ControllerInterface) {
	this.routerFunc = routerFunc
}

func (this *Router) AddController(path string, c ControllerInterface) {
	this.routerMap[path] = reflect.Indirect(reflect.ValueOf(c)).Type()
}

func (this *Router) AddAutoRouterController(path string, c ControllerInterface) {
	this.autoRouterMap[path] = reflect.Indirect(reflect.ValueOf(c)).Type()
}

func (this *Router) GetController(path string) ControllerInterface {
	var ci ControllerInterface = nil
	//router func should be first
	if this.routerFunc != nil {
		ci = this.routerFunc(path)
	}

	//normal router map is second
	if ci == nil {
		cit, _ := this.routerMap[path]
		if cit == nil {
			//auto router map is the last
			comps := toolkit.SplitString(path, "/")
			if len(comps) > 1 {
				cit, _ = this.autoRouterMap[comps[0]]
			}
		}

		if cit == nil {
			return nil
		} else {
			return reflect.New(cit).Interface().(ControllerInterface)
		}

	} else {
		return ci
	}
}
