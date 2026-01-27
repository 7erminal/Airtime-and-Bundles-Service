package helpers

import (
	"airtime_payment_service/models"
	"fmt"
	"strings"
	"time"

	"encoding/base64"

	"github.com/beego/beego/v2/core/logs"
)

func GetNetworkCode(networkName string, serviceType string) (resp string) {
	networkCode := networkName + "_" + serviceType

	return networkCode
}

func GetServiceId(network string) (string, error) {

	if networkService, err := models.GetNetworksByCode(network); err == nil {
		return networkService.NetworkReferenceId, nil
	}

	return "", nil
}

func Logger(logLevel string, requestId string, message string) {
	//do nothing for now
	// try to extract requestID from message if present as "requestID=" or "requestID:"

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	out := fmt.Sprintf("%s :: %s - %s ~ %s", timestamp, strings.ToUpper(logLevel), requestId, message)

	switch strings.ToLower(logLevel) {
	case "debug":
		logs.Debug(out)
	case "info":
		logs.Info(out)
	case "warn", "warning":
		logs.Warn(out)
	case "error":
		logs.Error(out)
	default:
		logs.Info(out)
	}
}

func ConvertToBase64(input string) string {
	encoded := ""
	// encoding logic here
	encoded = base64.StdEncoding.EncodeToString([]byte(input))
	return encoded
}

func ConvertFromBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
