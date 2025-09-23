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

type BilTransactionsData struct {
	TransactionId           int64
	TransactionRefNumber    string
	Service                 string
	Request                 int64
	TransactionBy           string
	Amount                  float64
	TransactingCurrency     string
	SourceChannel           string
	Source                  string
	Destination             string
	Charge                  float64
	BillerName              string
	NetworkName             string
	Commission              float64
	ExternalReferenceNumber string
	Status                  string
	DateCreated             string
	DateModified            string
	CreatedBy               int
	ModifiedBy              int
	Active                  int
}

type BilTransactionsResponse struct {
	StatusCode    string
	StatusMessage string
	Result        *BilTransactionsData
}

type BilTransactionsListResponse struct {
	StatusCode    string
	StatusMessage string
	Result        []*BilTransactionsData
	TotalCount    int64
}
