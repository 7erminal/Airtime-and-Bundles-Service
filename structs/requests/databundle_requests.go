package requests

type GetDataBundleRequest struct {
	PhoneNumber string `json:"phone_number" valid:"required~Phone number is required"`
	Network     string `json:"network" valid:"required~Network is required"`
}

type ThirdPartyGetDataBundlesRequest struct {
	Destination string `json:"destination" valid:"required~Phone number is required"`
}

type DataBundleRequest struct {
	Amount      float64 `json:"amount" valid:"required~Amount is required"`
	Network     string  `json:"network" valid:"required~Network is required"`
	Destination string  `json:"destination" valid:"required~Destination is required"`
	BundleId    string  `json:"bundle_id" valid:"required~Bundle ID is required"`
}

type BundleKeyRequest struct {
	Bundle string `json:"bundle" valid:"required~Bundle key is required"`
}

type DataBundleThirdPartyRequest struct {
	PhoneNumber     string
	Amount          float64
	Network         string
	Destination     string
	CallbackUrl     string
	ClientReference string
	ExtraData       BundleKeyRequest `json:"extra_data" valid:"required~Extra data is required"`
	BundleId        string
}
