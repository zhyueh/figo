package context

import (
	"github.com/zhyueh/figo/toolkit"
	"net/http"
	"sync"
)

var (
	//context data base on http request
	context_data_mutex sync.RWMutex
	context_data       = make(map[*http.Request]map[interface{}]interface{})
)

func RequestNum() int {
	return len(context_data)
}

func Clear(request *http.Request) {
	context_data_mutex.Lock()
	defer context_data_mutex.Unlock()

	delete(context_data, request)
}

func Delete(request *http.Request, key interface{}) {
	context_data_mutex.Lock()
	defer context_data_mutex.Unlock()

	if context_data[request] != nil {
		delete(context_data[request], key)
	}
}

func Set(request *http.Request, key, value interface{}) {
	context_data_mutex.Lock()
	defer context_data_mutex.Unlock()

	if context_data[request] == nil {
		context_data[request] = make(map[interface{}]interface{})
	}
	context_data[request][key] = value
}

func GetAll(request *http.Request) map[interface{}]interface{} {
	context_data_mutex.RLock()
	defer context_data_mutex.RUnlock()

	return context_data[request]
}

func GetOK(request *http.Request, key interface{}) (interface{}, bool) {
	context_data_mutex.RLock()
	defer context_data_mutex.RUnlock()

	request_context_data := GetAll(request)
	if request_context_data != nil {
		value, ok := request_context_data[key]
		return value, ok
	} else {
		return nil, false
	}
}

func Get(request *http.Request, key interface{}) interface{} {
	value, _ := GetOK(request, key)
	return value
}

func GetInt(request *http.Request, key interface{}) int {
	return toolkit.ConvertToInt(Get(request, key))
}

func GetInt64(request *http.Request, key interface{}) int64 {
	return toolkit.ConvertToInt64(Get(request, key))
}

func GetString(request *http.Request, key interface{}) string {
	return toolkit.ConvertToString(Get(request, key))
}
