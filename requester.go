package figo

import (
	"net/http"
)

type Requester struct {
	BaseHttpRequest *http.Request
}

func NewRequester(r *http.Request) *Requester {
	this := &Requester{}
	this.BaseHttpRequest = r

	this.init()
	return this
}

func (this *Requester) init() {
	//define your self base on your protocol

}
