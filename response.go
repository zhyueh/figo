package figo

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	BaseHttpResponse http.ResponseWriter
	extraHeader      map[string]string
	responseBody     []byte
	flushed          bool
}

func NewResponse(w http.ResponseWriter) *Response {
	this := &Response{}
	this.BaseHttpResponse = w
	this.extraHeader = make(map[string]string, 2)
	this.responseBody = make([]byte, 1024)
	this.flushed = false

	return this
}

func (this *Response) AddHeader(key string, value string) {
	this.extraHeader[key] = value
}

func (this *Response) WriteString(str string) {
	this.WriteByte([]byte(str))
}

func (this *Response) WriteJson(obj interface{}) error {
	jsonData, err := json.Marshal(obj)
	if err == nil {
		this.WriteByte(jsonData)
	}
	return err
}

func (this *Response) WriteByte(data []byte) {
	this.responseBody = data
}

// suppose to be execute in figo app
func (this *Response) Flush() error {
	if !this.flushed {
		//flush to base http response writer

		this.BaseHttpResponse.Write(this.responseBody)
		for k, v := range this.extraHeader {
			this.BaseHttpResponse.Header().Add(k, v)
		}

		this.flushed = true
		return nil
	}
	return errors.New("we can not flush twice")
}
