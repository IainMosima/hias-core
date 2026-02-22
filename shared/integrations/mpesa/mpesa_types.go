package mpesa

import "time"

type STKPushRequest struct {
	PhoneNumber string `json:"phone_number"`
	Amount      int64  `json:"amount"`
	AccountRef  string `json:"account_reference"`
	Description string `json:"description"`
}

type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

type B2CRequest struct {
	PhoneNumber string `json:"phone_number"`
	Amount      int64  `json:"amount"`
	Occasion    string `json:"occasion"`
	Remarks     string `json:"remarks"`
}

type B2CResponse struct {
	ConversationID          string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode            string `json:"ResponseCode"`
	ResponseDescription     string `json:"ResponseDescription"`
}

type CallbackData struct {
	MerchantRequestID string    `json:"MerchantRequestID"`
	CheckoutRequestID string    `json:"CheckoutRequestID"`
	ResultCode        int       `json:"ResultCode"`
	ResultDesc        string    `json:"ResultDesc"`
	Amount            int64     `json:"Amount,omitempty"`
	MpesaReceiptNo    string    `json:"MpesaReceiptNumber,omitempty"`
	TransactionDate   time.Time `json:"TransactionDate,omitempty"`
	PhoneNumber       string    `json:"PhoneNumber,omitempty"`
}

type StatusQueryRequest struct {
	CheckoutRequestID string `json:"checkout_request_id"`
}

type StatusQueryResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	ResultCode          string `json:"ResultCode"`
	ResultDesc          string `json:"ResultDesc"`
}

type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}
