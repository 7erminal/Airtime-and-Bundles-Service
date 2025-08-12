package responses

type AirtimeResponseResult struct {
	TransactionID     string
	PhoneNumber       string
	Amount            float64
	Network           string
	Destination       string
	TransactionStatus string
	TransactionDate   string
}

type AirtimeResponse struct {
	StatusCode    bool
	StatusMessage string
	Result        *AirtimeResponseResult
}

type ThirdPartyMeta struct {
	Commission string
}

type ThirdPartyDataResponse struct {
	ClientReference string
	Amount          float64
	TransactionId   string
	Meta            ThirdPartyMeta
}

type ThirdPartyAirtimeResponse struct {
	ResponseCode string
	Message      string
	Data         ThirdPartyDataResponse
}
