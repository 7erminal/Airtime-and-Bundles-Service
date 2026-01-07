package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/astaxie/beego/httplib"
	"github.com/beego/beego/v2/core/logs"
)

type Client struct {
	Request *Request
	Type_   string
}

func (c *Client) SendRequest() (*http.Response, error) {
	// TODO: BaseURLとPathのvalidation
	fullURL := c.Request.BaseURL + c.Request.Path
	logs.Info("Full third party URL: ", fullURL)
	beegoRequest := httplib.NewBeegoRequest(
		fullURL,
		c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	for key, value := range c.Request.FileField {
		logs.Info("File field seen...")
		logs.Info("Key is ", key, " and value is ", value)
		fileInfo, err := os.Stat(value)
		if err != nil {
			logs.Error("Error accessing file:", err)
			value = ""
		} else {
			logs.Info("File size:", fileInfo.Size())
		}
		beegoRequest.PostFile(key, value)
	}
	if c.Type_ == "params" {
		for key, value := range c.Request.Params {
			beegoRequest.Param(key, value)
		}
	} else if c.Type_ == "body" {
		reqText := []byte{}
		if c.Request.InterfaceParams == nil {
			logs.Info("Interface params are nil")
			beegoRequest.JSONBody(c.Request.Params)

			reqText, _ = json.Marshal(c.Request.Params)
		} else {
			logs.Info("Interface params are not nil")
			beegoRequest.JSONBody(c.Request.InterfaceParams)

			reqText, _ = json.Marshal(c.Request.InterfaceParams)
		}

		logs.Info("Request to third party is ", string(reqText))
	}

	beegoRequest.Debug(true)
	logs.Info("beego request is ", beegoRequest)

	res, err := beegoRequest.Response()
	if err != nil {
		return nil, err
	}
	return res, nil
}
