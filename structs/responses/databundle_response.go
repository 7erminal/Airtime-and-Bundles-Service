package responses

type DataBundleResponseResult struct {
	TransactionID     string  `json:"transaction_id"`
	PhoneNumber       string  `json:"phone_number"`
	Amount            float64 `json:"amount"`
	Network           string  `json:"network"`
	Destination       string  `json:"destination"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionDate   string  `json:"transaction_date"`
}
type DataBundleResponse struct {
	StatusCode    bool
	StatusMessage string
	Result        *DataBundleResponseResult
}
type ThirdPartyBundleMeta struct {
	Commission string `json:"commission"`
}
type ThirdPartyBundleDataResponse struct {
	ClientReference string         `json:"client_reference"`
	Amount          float64        `json:"amount"`
	TransactionId   string         `json:"transaction_id"`
	Meta            ThirdPartyMeta `json:"meta"`
}
type ThirdPartyDataBundleResponse struct {
	ResponseCode string
	Message      string
	Data         ThirdPartyBundleDataResponse
}

type DataBundleData struct {
	Display string
	Value   string
	Amount  float64
}

type ThirdPartyDataBundlesResponse struct {
	ResponseCode string
	Message      string
	Data         []DataBundleData
}

type GetBundlesResponse struct {
	StatusCode    bool
	StatusMessage string
	Result        []DataBundleData
}
