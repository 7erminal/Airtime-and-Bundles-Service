package requests

type AirtimeRequest struct {
	Amount      float64 `json:"amount" valid:"required~Amount is required"`
	Network     string  `json:"network" valid:"required~Network is required"`
	Destination string  `json:"destination" valid:"required~Destination is required"`
}

type AirtimeThirdPartyRequest struct {
	PhoneNumber   string  `json:"phone_number" valid:"required~Phone number is required"`
	Amount        float64 `json:"amount" valid:"required~Amount is required"`
	Destination   string  `json:"destination" valid:"required~Destination is required"`
	TransactionId string  `json:"transaction_id" valid:"required~Transaction ID is required"`
	Network       string  `json:"network" valid:"required~Network is required"` // Assuming this is the network name
	CallbackUrl   string  `json:"callback_url" valid:"optional"`                // Optional field for callback URL
}
