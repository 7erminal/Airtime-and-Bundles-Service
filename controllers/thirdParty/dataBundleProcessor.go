package thirdparty

import (
	"airtime_payment_service/api"
	"airtime_payment_service/helpers"
	"airtime_payment_service/structs/requests"
	"airtime_payment_service/structs/responses"
	"encoding/json"
	"io"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func ProcessDataBundlePurchase(c *beego.Controller, req requests.DataBundleThirdPartyRequest) (responses.ThirdPartyDataBundleResponse, error) {
	host, _ := beego.AppConfig.String("thirdPartyBaseUrl")
	prepaidId, _ := beego.AppConfig.String("hubtelPrepaidDepositID")
	authorizationKey, _ := beego.AppConfig.String("authorizationKey")

	logs.Info("Sending phone number ", req.PhoneNumber)
	logs.Info("Using network ", req.Network)

	serviceId, _ := helpers.GetServiceId(req.Network)

	reqText, _ := json.Marshal(req)

	logs.Info("Request to process data bundle purchase: ", string(reqText))

	request := api.NewRequest(
		host,
		"/"+prepaidId+"/"+serviceId,
		api.POST)
	request.HeaderField["Authorization"] = "Basic " + authorizationKey
	request.InterfaceParams["Destination"] = req.Destination
	request.InterfaceParams["Amount"] = req.Amount
	request.InterfaceParams["CallbackUrl"] = req.CallbackUrl
	request.InterfaceParams["ClientReference"] = req.ClientReference
	request.InterfaceParams["ExtraData"] = req.ExtraData
	// request.InterfaceParams["BundleId"] = req.BundleId

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "body",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	logs.Info("Raw response received is ", res)
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ThirdPartyDataBundleResponse
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetDataBundles(c *beego.Controller, networkCode string, destinationNumber string) (responses.ThirdPartyDataBundlesResponse, error) {
	host, _ := beego.AppConfig.String("thirdPartyBaseUrl")
	prepaidId, _ := beego.AppConfig.String("hubtelPrepaidDepositID")
	authorizationKey, _ := beego.AppConfig.String("authorizationKey")

	logs.Info("Authorization key is " + authorizationKey)

	// serviceId, _ := helpers.GetServiceId(networkCode)

	request := api.NewRequest(
		host,
		"/"+prepaidId+"/"+networkCode+"?destination="+destinationNumber,
		api.GET)
	request.HeaderField["Authorization"] = "Basic " + authorizationKey

	// request.InterfaceParams["BundleId"] = req.BundleId

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	// request.Params["destination"] = destinationNumber
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	logs.Info("Raw response received is ", res)
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ThirdPartyDataBundlesResponse
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}
