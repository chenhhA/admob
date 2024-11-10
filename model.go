package admob

import (
	"net/http"
	"time"
)

type Config struct {
	// PublicKeyCacheDuration public keys are regularly rotated and should not be cached for longer than 24 hours.
	PublicKeyCacheDuration time.Duration
	HttpClient             *http.Client
}

type CallbackParam struct {
	AdNetwork     string
	AdUnit        string
	CustomData    string
	KeyID         int
	RewardAmount  string
	RewardItem    string
	Signature     string
	Timestamp     string
	TransactionID string
	UserID        string
}

type adModResponse struct {
	Keys []*publicKey `json:"keys"`
}

type publicKey struct {
	KeyId  int    `json:"keyId"`
	Pem    string `json:"pem"`
	Base64 string `json:"base64"`
}
