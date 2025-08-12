package controllers

import (
	thirdparty "airtime_payment_service/controllers/thirdParty"
	"airtime_payment_service/helpers"
	"airtime_payment_service/models"
	"airtime_payment_service/structs/requests"
	"airtime_payment_service/structs/responses"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// RequestController operations for Request
type RequestController struct {
	beego.Controller
}

// URLMapping ...
func (c *RequestController) URLMapping() {
	c.Mapping("BuyAirtime", c.BuyAirtime)
	c.Mapping("BuyDataBundle", c.BuyDataBundle)
	c.Mapping("GetBundles", c.GetBundles)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Request
// @Param	PhoneNumber		header 	string true		"header for Customer's phone number"
// @Param	SourceSystem		header 	string true		"header for Source system"
// @Param	body		body 	requests.AirtimeRequest	true		"body for Request content"
// @Success 201 {int} models.Request
// @Failure 403 body is empty
// @router /buy-airtime [post]
func (c *RequestController) BuyAirtime() {
	var req requests.AirtimeRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	// Validate the request
	phoneNumber := c.Ctx.Input.Header("PhoneNumber")
	sourceSystem := c.Ctx.Input.Header("SourceSystem")

	responseCode := false
	responseMessage := "Request not processed"

	statusCode := "PENDING" // Assuming 5002 is the status code for "Request Pending"

	reqText, err := json.Marshal(req)
	if err != nil {
		c.Data["json"] = "Invalid request format"
		c.ServeJSON()
		return
	}

	status, err := models.GetStatus_codesByCode(statusCode)
	if err == nil {
		// Get customer by ID
		logs.Info("About to get customer by phone number: ", phoneNumber)
		if cust, err := models.GetCustomerByPhoneNumber(phoneNumber); err == nil {
			// Restructure the request to match the model
			serviceCode := "AIRTIME"
			if service, err := models.GetServicesByCode(serviceCode); err == nil {
				v := models.Request{
					CustId:          cust,
					Request:         string(reqText),
					RequestType:     service.ServiceName,
					RequestStatus:   status.StatusDescription,
					RequestAmount:   req.Amount,
					RequestResponse: "",
					RequestDate:     time.Now(),
					DateCreated:     time.Now(),
					DateModified:    time.Now(),
				}
				if _, err := models.AddRequest(&v); err == nil {
					// Create a transaction record
					currentTime := time.Now().Unix()
					transaction := models.Bil_transactions{
						TransactionRefNumber: "TRX-" + strconv.FormatInt(currentTime, 10) + strconv.FormatInt(v.RequestId, 10),
						Service:              service, // Assuming service ID is 1 for airtime
						Request:              &v,
						TransactionBy:        cust,
						Amount:               req.Amount,
						TransactingCurrency:  "GHC", // Assuming USD for simplicity
						SourceChannel:        sourceSystem,
						Source:               phoneNumber,
						Destination:          req.Destination,
						Charge:               0.0,    // Assuming no charge for simplicity
						Status:               status, // Assuming 1 means successful
						DateCreated:          time.Now(),
						DateModified:         time.Now(),
						CreatedBy:            1,
						ModifiedBy:           1,
						Active:               1, // Assuming active status
					}
					if _, err := models.AddBil_transactions(&transaction); err == nil {
						// Go to fulfillment
						callbackurl := ""
						if cbr, err := models.GetApplication_propertyByCode("AIRTIME_CALLBACK_URL"); err == nil {
							logs.Info("Property data is ", cbr)
							logs.Info("Callback url is ", cbr.PropertyValue)
							callbackurl = cbr.PropertyValue
						} else {
							logs.Error("Failed to get callback URL: %v", err)
						}

						networkCode := helpers.GetNetworkCode(req.Network, service.ServiceCode)

						billerCode := "AIRTIME"
						biller, err := models.GetBillerByCode(billerCode)

						if err == nil {
							// Formulate the request to send to the third-party service
							tReq := requests.AirtimeThirdPartyRequest{
								PhoneNumber:   phoneNumber,
								Amount:        req.Amount,
								Destination:   req.Destination,
								Network:       networkCode,                      // Assuming service name is used as network
								TransactionId: transaction.TransactionRefNumber, // Use the request ID as the transaction ID
								CallbackUrl:   callbackurl,                      // Optional field for callback URL
							}

							// Insert in INS Transactions table
							reqText, err := json.Marshal(tReq)
							if err != nil {
								logs.Error("Failed to marshal request text: %v", err)
								// c.Data["json"] = "Invalid request format"
								// c.ServeJSON()
								// return
							}

							insTransaction := models.Bil_ins_transactions{
								BilTransactionId:       &transaction,
								Amount:                 req.Amount,
								Biller:                 biller,
								SenderAccountNumber:    phoneNumber,
								RecipientAccountNumber: req.Destination,
								Network:                billerCode,
								Request:                string(reqText),
								DateCreated:            time.Now(),
								DateModified:           time.Now(),
								CreatedBy:              1,
								ModifiedBy:             1,
								Active:                 1,
							}

							if _, err := models.AddBil_ins_transactions(&insTransaction); err != nil {
								logs.Error("Failed to create INS transaction record: %v", err)
								responseCode = false
								responseMessage = "Failed to create INS transaction record"
								// resp := responses.ThirdPartyBillPaymentApiResponse{
								// 	StatusCode:    responseCode,
								// 	StatusMessage: responseMessage,
								// 	Result:        nil,
								// }
								// c.Data["json"] = resp
								// c.ServeJSON()
								// return
							}

							// Call the third-party service to process the request
							logs.Info("Processing airtime request with third-party service: ", tReq)
							if thirdPartyResponse, err := thirdparty.ProcessAirtime(&c.Controller, tReq); err == nil {

								if thirdPartyResponse.ResponseCode == "0001" {
									// Transaction is pending
									// Update the transaction status to pending
									responseCode = true
									responseMessage = "Request is being processed"
									if status, err := models.GetStatus_codesByCode("PENDING"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "PENDING:: Failed to update transaction status"
										} else {
											responseCode = true
											responseMessage = "Request is being processed"
										}
									} else {
										logs.Error("Failed to get status for pending transaction: %v", err)
										responseCode = false
										responseMessage = "PENDING: Failed to get status for pending transaction"
									}
								} else if thirdPartyResponse.ResponseCode == "0000" {
									// Transaction is successful
									// Update the transaction status to successful
									responseCode = true
									responseMessage = "Request is successful"
									if status, err := models.GetStatus_codesByCode("SUCCESS"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "SUCCESS:: Failed to update transaction status"
										} else {
											// Prepare the response
											logs.Info("Transaction successful: ", transaction)
											responseCode = true
											responseMessage = "Transaction successful"
										}
									} else {
										logs.Error("Failed to get status for successful transaction: %v", err)
										responseCode = false
										responseMessage = "SUCCESS:: Failed to get status for successful transaction"
									}
								} else {
									// Transaction failed
									// Update the transaction status to failed
									responseCode = false
									responseMessage = "Transaction failed"
									if status, err := models.GetStatus_codesByCode("FAILED"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "FAILED:: Failed to update transaction status"
										}
									} else {
										logs.Error("Failed to get status for failed transaction: %v", err)
										responseCode = false
										responseMessage = "FAILED:: Failed to get status for failed transaction"
									}
								}
								// Update requests response
								resText, err := json.Marshal(thirdPartyResponse)
								if err != nil {
									logs.Error("Failed to marshal response text: %v", err)
									// c.Data["json"] = "Invalid request format"
									// c.ServeJSON()
									// return
								}
								v.RequestResponse = string(resText)
								v.DateModified = time.Now()
								if err := models.UpdateRequestById(&v); err != nil {
									logs.Error("Failed to update request response: %v", err)
									responseCode = true
									responseMessage = "Success response:: Failed to update request response"
								} else {
									logs.Info("Request response updated successfully")
								}
								c.Ctx.Output.SetStatus(200)
								// Prepare the response

								// Create the response object
								airtimeResponse := responses.AirtimeResponseResult{
									PhoneNumber:       cust.PhoneNumber,
									Amount:            req.Amount,
									Network:           req.Network,
									Destination:       req.Destination,
									TransactionStatus: status.StatusDescription,
									TransactionDate:   v.RequestDate.Format(time.RFC3339),
								}
								response := responses.AirtimeResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        &airtimeResponse,
								}
								c.Data["json"] = response
							} else {
								logs.Error("Failed to process third-party request: %v", err)
								responseCode = false
								responseMessage = "Failed to process third-party request"
								resp := responses.AirtimeResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        nil,
								}
								c.Data["json"] = resp
							}
						} else {
							logs.Error("Failed to get biller: %v", err)
							responseCode = false
							responseMessage = "Failed to get biller"
							resp := responses.AirtimeResponse{
								StatusCode:    responseCode,
								StatusMessage: responseMessage,
								Result:        nil,
							}
							c.Data["json"] = resp
						}
					} else {
						logs.Error("Failed to create transaction record: %v", err)
						responseCode = false
						responseMessage = "Failed to create transaction record"
						resp := responses.AirtimeResponse{
							StatusCode:    responseCode,
							StatusMessage: responseMessage,
							Result:        nil,
						}
						c.Data["json"] = resp
					}
				} else {
					logs.Error("Failed to create request record: %v", err)
					responseCode = false
					responseMessage = "Failed to create transaction record"
					resp := responses.AirtimeResponse{
						StatusCode:    responseCode,
						StatusMessage: responseMessage,
						Result:        nil,
					}
					c.Data["json"] = resp
				}
			} else {
				logs.Error("Service not found: %v", err)
				responseCode = false
				responseMessage = "Failed to create transaction record"
				resp := responses.AirtimeResponse{
					StatusCode:    responseCode,
					StatusMessage: responseMessage,
					Result:        nil,
				}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("Customer not found: %v", err)
			responseCode = false
			responseMessage = "Failed to create transaction record"
			resp := responses.AirtimeResponse{
				StatusCode:    responseCode,
				StatusMessage: responseMessage,
				Result:        nil,
			}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Status not found: %v", err)
		responseCode = false
		responseMessage = "Failed to create transaction record"
		resp := responses.AirtimeResponse{
			StatusCode:    responseCode,
			StatusMessage: responseMessage,
			Result:        nil,
		}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// DataBundles ...
// @Title Buy Data
// @Description create Request
// @Param	PhoneNumber		header 	string true		"header for Customer's phone number"
// @Param	SourceSystem		header 	string true		"header for Source system"
// @Param	body		body 	requests.DataBundleRequest	true		"body for Request content"
// @Success 201 {int} models.Request
// @Failure 403 body is empty
// @router /buy-bundle [post]
func (c *RequestController) BuyDataBundle() {
	var req requests.DataBundleRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	// Validate the request

	// authorization := ctx.Input.Header("Authorization")
	phoneNumber := c.Ctx.Input.Header("PhoneNumber")
	sourceSystem := c.Ctx.Input.Header("SourceSystem")

	responseCode := false
	responseMessage := "Request not processed"

	statusCode := "PENDING" // Assuming 5002 is the status code for "Request Pending"

	reqText, err := json.Marshal(req)
	if err != nil {
		c.Data["json"] = "Invalid request format"
		c.ServeJSON()
		return
	}

	status, err := models.GetStatus_codesByCode(statusCode)
	if err == nil {
		// Get customer by ID
		if cust, err := models.GetCustomerByPhoneNumber(phoneNumber); err == nil {
			// Restructure the request to match the model
			serviceCode := "DATA_BUNDLE"
			if service, err := models.GetServicesByCode(serviceCode); err == nil {
				v := models.Request{
					CustId:          cust,
					Request:         string(reqText),
					RequestType:     service.ServiceName,
					RequestStatus:   status.StatusDescription,
					RequestAmount:   req.Amount,
					RequestResponse: "",
					RequestDate:     time.Now(),
					DateCreated:     time.Now(),
					DateModified:    time.Now(),
				}
				if _, err := models.AddRequest(&v); err == nil {
					// Create a transaction record
					transaction := models.Bil_transactions{
						TransactionRefNumber: "TRX-" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(v.RequestId, 10),
						Service:              service, // Assuming service ID is 1 for airtime
						Request:              &v,
						TransactionBy:        cust,
						Amount:               req.Amount,
						TransactingCurrency:  "GHC", // Assuming USD for simplicity
						SourceChannel:        sourceSystem,
						Source:               phoneNumber,
						Destination:          req.Destination,
						Charge:               0.0,    // Assuming no charge for simplicity
						Status:               status, // Assuming 1 means successful
						DateCreated:          time.Now(),
						DateModified:         time.Now(),
						CreatedBy:            1,
						ModifiedBy:           1,
						Active:               1, // Assuming active status
					}
					if _, err := models.AddBil_transactions(&transaction); err == nil {
						// Go to fulfillment
						// Formulate the request to send to the third-party service
						selectedBundle := requests.BundleKeyRequest{
							Bundle: req.BundleId,
						}

						callbackurl := ""
						if cbr, err := models.GetApplication_propertyByCode("DATA_BUNDLE_CALLBACK_URL"); err == nil {
							callbackurl = cbr.PropertyValue
						} else {
							logs.Error("Failed to get callback URL: %v", err)
						}

						logs.Info("About to get network code for service: ", service.ServiceCode, " and network: ", req.Network)

						networkCode := helpers.GetNetworkCode(req.Network, service.ServiceCode)

						billerCode := "DATA_BUNDLE"
						biller, err := models.GetBillerByCode(billerCode)

						if err == nil {

							tReq := requests.DataBundleThirdPartyRequest{
								PhoneNumber:     phoneNumber,
								Amount:          req.Amount,
								Destination:     req.Destination,
								Network:         networkCode,                      // Assuming service name is used as network
								ClientReference: transaction.TransactionRefNumber, // Use the request ID as the transaction ID
								CallbackUrl:     callbackurl,                      // Optional field for callback URL
								ExtraData:       selectedBundle,                   // Assuming this is the bundle key request
								BundleId:        req.BundleId,                     // Assuming this is the bundle ID

							}

							// Insert in INS Transactions table
							reqText, err := json.Marshal(tReq)
							if err != nil {
								logs.Error("Failed to marshal request text: %v", err)
								// c.Data["json"] = "Invalid request format"
								// c.ServeJSON()
								// return
							}

							insTransaction := models.Bil_ins_transactions{
								BilTransactionId:       &transaction,
								Amount:                 req.Amount,
								Biller:                 biller,
								SenderAccountNumber:    phoneNumber,
								RecipientAccountNumber: req.Destination,
								Network:                billerCode,
								Request:                string(reqText),
								DateCreated:            time.Now(),
								DateModified:           time.Now(),
								CreatedBy:              1,
								ModifiedBy:             1,
								Active:                 1,
							}

							if _, err := models.AddBil_ins_transactions(&insTransaction); err != nil {
								logs.Error("Failed to create INS transaction record: %v", err)
								responseCode = false
								responseMessage = "Failed to create INS transaction record"
								// resp := responses.ThirdPartyBillPaymentApiResponse{
								// 	StatusCode:    responseCode,
								// 	StatusMessage: responseMessage,
								// 	Result:        nil,
								// }
								// c.Data["json"] = resp
								// c.ServeJSON()
								// return
							}

							// Call the third-party service to process the request
							logs.Info("Processing bundle request with third-party service: ", tReq)
							if thirdPartyResponse, err := thirdparty.ProcessDataBundlePurchase(&c.Controller, tReq); err == nil {
								logs.Info("Third-party response received: ", thirdPartyResponse)

								if thirdPartyResponse.ResponseCode == "0001" {
									// Transaction is pending
									// Update the transaction status to pending
									logs.Info("Transaction is pending, updating status to PENDING")
									responseCode = true
									responseMessage = "Request is being processed"
									if status, err := models.GetStatus_codesByCode("PENDING"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "PENDING:: Failed to update transaction status"
										} else {
											logs.Info("Transaction status updated to pending")
											responseCode = true
											responseMessage = "Request is being processed"
										}
									} else {
										logs.Error("Failed to get status for pending transaction: %v", err)
										responseCode = false
										responseMessage = "PENDING: Failed to get status for pending transaction"
									}
								} else if thirdPartyResponse.ResponseCode == "0000" {
									// Transaction is successful
									// Update the transaction status to successful
									responseCode = true
									responseMessage = "Request is successful"
									if status, err := models.GetStatus_codesByCode("SUCCESS"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "SUCCESS:: Failed to update transaction status"
										} else {
											// Prepare the response
											logs.Info("Transaction successful: ", transaction)
											responseCode = true
											responseMessage = "Transaction successful"
										}
									} else {
										logs.Error("Failed to get status for successful transaction: %v", err)
										responseCode = false
										responseMessage = "SUCCESS:: Failed to get status for successful transaction"
									}
								} else {
									// Transaction failed
									// Update the transaction status to failed
									responseCode = false
									responseMessage = "Transaction failed"
									if status, err := models.GetStatus_codesByCode("FAILED"); err == nil {
										transaction.Status = status
										if err := models.UpdateBil_transactionsById(&transaction); err != nil {
											logs.Error("Failed to update transaction status: %v", err)
											responseCode = false
											responseMessage = "FAILED:: Failed to update transaction status"
										}
									} else {
										logs.Error("Failed to get status for failed transaction: %v", err)
										responseCode = false
										responseMessage = "FAILED:: Failed to get status for failed transaction"
									}
								}

								resText, err := json.Marshal(thirdPartyResponse)
								if err != nil {
									logs.Error("Failed to marshal response text: %v", err)
									// c.Data["json"] = "Invalid request format"
									// c.ServeJSON()
									// return
								}
								v.RequestResponse = string(resText)
								v.DateModified = time.Now()
								if err := models.UpdateRequestById(&v); err != nil {
									logs.Error("Failed to update request response: %v", err)
									responseCode = false
									responseMessage = "Success response:: Failed to update request response"
								} else {
									logs.Info("Request response updated successfully")
									logs.Info("Response code: ", responseCode)
								}

								c.Ctx.Output.SetStatus(200)
								// Prepare the response

								// Create the response object
								bundleResponse := responses.DataBundleResponseResult{
									PhoneNumber:       cust.PhoneNumber,
									Amount:            req.Amount,
									Network:           req.Network,
									Destination:       req.Destination,
									TransactionStatus: status.StatusDescription,
									TransactionDate:   v.RequestDate.Format(time.RFC3339),
								}

								logs.Info("Bundle response created: ", bundleResponse)
								logs.Info("Response code: ", responseCode)
								logs.Info("Response message: ", responseMessage)
								response := responses.DataBundleResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        &bundleResponse,
								}
								c.Data["json"] = response
							} else {
								logs.Error("Failed to process third-party request: %v", err)
								responseCode = false
								responseMessage = "Failed to process third-party request"
								resp := responses.DataBundleResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        nil,
								}
								c.Data["json"] = resp
							}
						} else {
							logs.Error("Failed to create transaction record: %v", err)
							responseCode = false
							responseMessage = "Failed to create transaction record"
							resp := responses.DataBundleResponse{
								StatusCode:    responseCode,
								StatusMessage: responseMessage,
								Result:        nil,
							}
							c.Data["json"] = resp
						}
					} else {
						logs.Error("Failed to create transaction record: %v", err)
						responseCode = false
						responseMessage = "Failed to create transaction record"
						resp := responses.DataBundleResponse{
							StatusCode:    responseCode,
							StatusMessage: responseMessage,
							Result:        nil,
						}
						c.Data["json"] = resp
					}
				} else {
					logs.Error("Failed to create request record: %v", err)
					responseCode = false
					responseMessage = "Failed to create transaction record"
					resp := responses.DataBundleResponse{
						StatusCode:    responseCode,
						StatusMessage: responseMessage,
						Result:        nil,
					}
					c.Data["json"] = resp
				}
			} else {
				logs.Error("Service not found: %v", err)
				responseCode = false
				responseMessage = "Failed to create transaction record"
				resp := responses.DataBundleResponse{
					StatusCode:    responseCode,
					StatusMessage: responseMessage,
					Result:        nil,
				}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("Customer not found: %v", err)
			responseCode = false
			responseMessage = "Failed to create transaction record"
			resp := responses.DataBundleResponse{
				StatusCode:    responseCode,
				StatusMessage: responseMessage,
				Result:        nil,
			}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Status not found: %v", err)
		responseCode = false
		responseMessage = "Failed to create transaction record"
		resp := responses.DataBundleResponse{
			StatusCode:    responseCode,
			StatusMessage: responseMessage,
			Result:        nil,
		}
		c.Data["json"] = resp
	}
	logs.Info("Response sent: ", c.Data["json"])
	c.ServeJSON()
}

// GetBundles ...
// @Title Get Bundles
// @Description get data Bundles
// @Param	PhoneNumber		header 	string true		"header for Customer's phone number"
// @Param	SourceSystem		header 	string true		"header for Source system"
// @Param	networkId		path 	string	true		"The key for staticblock"
// @Param	destinationPhoneNumber		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Request
// @Failure 403 :network is empty
// @router /bundles/:networkId/:destinationPhoneNumber [get]
func (c *RequestController) GetBundles() {
	// idStr := c.Ctx.Input.Param(":id")

	// var sourceSystemQuery string
	phoneNumber := c.Ctx.Input.Header("PhoneNumber")
	sourceSystem := c.Ctx.Input.Header("SourceSystem")
	// id, _ := strconv.ParseInt(idStr, 0, 64)
	// var network string
	logs.Info("Request received")
	logs.Info("Customer's Phone number is ", phoneNumber)
	logs.Info("Source system is ", sourceSystem)

	// if sourceSystemQuery != "" {
	// 	sourceSystemQuery = sourceSystemQuery + ","
	// }
	// sourceSystemQuery = sourceSystemQuery + "Operator__operator_name:" + sourceSystem

	// logs.Info("About to query with source system ", sourceSystemQuery)
	// if v := sourceSystemQuery; v != "" {
	// 	for _, cond := range strings.Split(v, ",") {
	// 		kv := strings.SplitN(cond, ":", 2)
	// 		if len(kv) != 2 {
	// 			logs.Error("Invalid query key/value pair")
	// 			c.Data["json"] = errors.New("error: invalid query key/value pair")
	// 			c.ServeJSON()
	// 			return
	// 		}
	// 		k, v := kv[0], kv[1]
	// 		query[k] = v
	// 	}
	// }

	responseCode := false
	responseMessage := "Request not processed"

	statusCode := "PENDING"

	logs.Info("About to go fetch networks")

	networkIdStr := c.Ctx.Input.Param(":networkId")
	networkId, err := strconv.ParseInt(networkIdStr, 10, 64)
	if err != nil {
		logs.Error("Error converting network ID ", err)
	}
	destinationPhoneNumber := c.Ctx.Input.Param(":destinationPhoneNumber")
	logs.Info("Network ID is ", networkIdStr)
	logs.Info("Destination phone number is ", destinationPhoneNumber)
	proceed := false
	service := models.Services{}
	serviceName := "DATA_BUNDLE"
	if service_, err := models.GetServicesByCode(serviceName); err == nil {
		service = *service_
		proceed = true
	} else {
		logs.Error("Failed to get service: %v", err)
		responseCode = false
		responseMessage = "Failed to get service"
		c.Data["json"] = responses.GetBundlesResponse{
			StatusCode:    responseCode,
			StatusMessage: responseMessage,
			Result:        nil,
		}

	}

	network := models.Networks{}
	if proceed {
		proceed = false
		if network_, err := models.GetNetworksById(networkId); err != nil {
			logs.Error("An error occurred getting networks... ", err.Error())
			c.Data["json"] = err.Error()
			logs.Info("About to check for the code:: ", networkIdStr+"_"+service.ServiceCode)
			if network_, err := models.GetNetworksByCode(networkIdStr + "_" + service.ServiceCode); err != nil {
				logs.Error("An error occurred getting networks... ", err.Error())
				c.Data["json"] = err.Error()
			} else {
				logs.Info("Network found by code: ", network_)
				network = *network_
				proceed = true

			}
		} else {
			network = *network_
			proceed = true
		}
	}

	if proceed {
		logs.Info("Getting customer by phone number ", phoneNumber)
		if cust, err := models.GetCustomerByPhoneNumber(phoneNumber); err == nil {
			logs.Info("Checking status code")
			if status, err := models.GetStatus_codesByCode(statusCode); err == nil {
				serviceName := "FETCH"
				if service, err := models.GetServicesByName(serviceName); err == nil {
					v := models.Request{
						CustId:          cust,
						Request:         "/" + networkIdStr + "/" + destinationPhoneNumber,
						RequestType:     service.ServiceName,
						RequestStatus:   status.StatusDescription,
						RequestAmount:   0.0,
						RequestResponse: "",
						RequestDate:     time.Now(),
						DateCreated:     time.Now(),
						DateModified:    time.Now(),
					}
					if _, err := models.AddRequest(&v); err == nil {
						// Assuming the networkId is passed as a query parameter

						// networkId, err := strconv.ParseInt(networkIdStr, 10, 64)
						// if err != nil {
						// 	c.Data["json"] = "Invalid network ID"
						// 	c.ServeJSON()
						// 	return
						// }
						// thirdPartyNetworkId := ""
						// for _, n := range networks {
						// 	network, ok := n.(models.Networks)
						// 	if !ok {
						// 		continue
						// 	}
						// 	if network.NetworkId == networkId {
						// 		// Get bundles for the network
						// 		thirdPartyNetworkId = network.NetworkReferenceId
						// 		break
						// 	}
						// }
						thirdPartyNetworkId := network.NetworkReferenceId

						// call the third-party service to get bundles
						if thirdPartyNetworkId == "" {
							c.Data["json"] = "Network not found"
							c.ServeJSON()
							return
						}
						// networkCode := helpers.GetNetworkCode(req.Network, service.ServiceCode)
						logs.Info("Fetching bundles for network ID: ", thirdPartyNetworkId)
						bundles, err := thirdparty.GetDataBundles(&c.Controller, thirdPartyNetworkId, destinationPhoneNumber)

						if err == nil {
							if bundles.ResponseCode == "0000" {
								// Transaction is successful
								// Update the transaction status to successful
								responseCode = true
								responseMessage = "Request is successful"
								// Prepare the response
								logs.Info("Bundles fetched successfully: ", bundles)
								resp := responses.GetBundlesResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        bundles.Data,
								}
								c.Data["json"] = resp
							} else {
								// Transaction failed
								// Update the transaction status to failed
								responseCode = false
								responseMessage = "Transaction failed"
								logs.Error("Failed to fetch bundles: ", bundles.Message)
								// c.Data["json"] = bundles
								resp := responses.GetBundlesResponse{
									StatusCode:    responseCode,
									StatusMessage: responseMessage,
									Result:        nil,
								}
								c.Data["json"] = resp
							}
						} else {
							logs.Error("Failed to fetch bundles: %v", err)
							responseCode = false
							responseMessage = "Failed to fetch bundles"
							c.Data["json"] = responses.GetBundlesResponse{
								StatusCode:    responseCode,
								StatusMessage: responseMessage,
								Result:        nil,
							}
						}
					} else {
						logs.Error("Adding request failed: %v", err)
						responseCode = false
						responseMessage = "Failed to add request"
						resp := responses.GetBundlesResponse{
							StatusCode:    responseCode,
							StatusMessage: responseMessage,
							Result:        nil,
						}
						c.Data["json"] = resp
					}
				} else {
					logs.Error("Failed to get service: %v", err)
					responseCode = false
					responseMessage = "Failed to get service"
					c.Data["json"] = responses.GetBundlesResponse{
						StatusCode:    responseCode,
						StatusMessage: responseMessage,
						Result:        nil,
					}

				}
			} else {
				logs.Error("Failed to get customer by phone number: %v", err)
				responseCode = false
				responseMessage = "Failed to create transaction record"
				resp := responses.GetBundlesResponse{
					StatusCode:    responseCode,
					StatusMessage: responseMessage,
					Result:        nil,
				}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("Customer not found: %v", err)
			responseCode = false
			responseMessage = "Failed to create transaction record"
			resp := responses.GetBundlesResponse{
				StatusCode:    responseCode,
				StatusMessage: responseMessage,
				Result:        nil,
			}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Failed to get networks: %v", err)
		responseCode = false
		responseMessage = "Failed to get networks"
		c.Data["json"] = responses.GetBundlesResponse{
			StatusCode:    responseCode,
			StatusMessage: responseMessage,
			Result:        nil,
		}
	}

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Request by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Request
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RequestController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetRequestById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Request
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Request
// @Failure 403
// @router / [get]
func (c *RequestController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllRequest(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Request
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Request	true		"body for Request content"
// @Success 200 {object} models.Request
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RequestController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Request{RequestId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateRequestById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Request
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RequestController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteRequest(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
