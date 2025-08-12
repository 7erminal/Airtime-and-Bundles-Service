package helpers

import "airtime_payment_service/models"

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
