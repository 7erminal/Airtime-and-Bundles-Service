package thirdparty

import (
	"airtime_payment_service/api"
	"airtime_payment_service/helpers"
	"airtime_payment_service/structs/requests"
	"airtime_payment_service/structs/responses"
	"encoding/json"
	"io"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// containsIgnoreCase checks if substr is in s, case-insensitive.
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func ProcessAirtime(c *beego.Controller, req requests.AirtimeThirdPartyRequest) (responses.ThirdPartyAirtimeResponse, error) {
	host, _ := beego.AppConfig.String("thirdPartyBaseUrl")
	prepaidId, _ := beego.AppConfig.String("hubtelPrepaidDepositID")
	authorizationKey, _ := beego.AppConfig.String("authorizationKey")

	logs.Info("Sending phone number ", req.PhoneNumber)
	logs.Info("Network is ", req.Network)
	logs.Info("Callback URL is ", req.CallbackUrl)
	logs.Info("Amount is ", req.Amount)
	logs.Info("Destination is ", req.Destination)

	serviceId, _ := helpers.GetServiceId(req.Network)

	reqText, _ := json.Marshal(req)

	logs.Info("Request to process airtime purchase: ", string(reqText))

	request := api.NewRequest(
		host,
		"/"+prepaidId+"/"+serviceId,
		api.POST)
	request.HeaderField["Authorization"] = "Basic " + authorizationKey
	request.InterfaceParams["Destination"] = req.Destination
	request.InterfaceParams["Amount"] = req.Amount
	request.InterfaceParams["CallbackUrl"] = req.CallbackUrl
	request.InterfaceParams["ClientReference"] = req.TransactionId

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
	var data responses.ThirdPartyAirtimeResponse
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}
