package admob

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/patrickmn/go-cache"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Verifier struct {
	client *http.Client
	cache  *cache.Cache
	mutex  *sync.Mutex
}

func NewVerifier() *Verifier {
	return NewVerifyWithConfig(&Config{
		PublicKeyCacheDuration: time.Hour * 10,
		HttpClient:             http.DefaultClient,
	})
}

func NewVerifyWithConfig(config *Config) *Verifier {
	return &Verifier{
		client: config.HttpClient,
		cache:  cache.New(config.PublicKeyCacheDuration, time.Hour),
		mutex:  &sync.Mutex{},
	}
}

func (v *Verifier) Verify(cbUrl *url.URL) (*CallbackParam, error) {
	// escape query & parse
	rawQuery, err := url.QueryUnescape(cbUrl.RawQuery)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, err
	}
	callBackParam := &CallbackParam{
		AdNetwork:     query.Get("ad_network"),
		AdUnit:        query.Get("ad_unit"),
		RewardAmount:  query.Get("reward_amount"),
		CustomData:    query.Get("custom_data"),
		RewardItem:    query.Get("reward_item"),
		Signature:     query.Get("signature"),
		Timestamp:     query.Get("timestamp"),
		TransactionID: query.Get("transaction_id"),
		UserID:        query.Get("user_id"),
	}
	signatureIdx := strings.Index(rawQuery, "&signature")
	if signatureIdx == -1 {
		return nil, ErrInvalidSignature
	}
	msgHash := hash(rawQuery[0:signatureIdx])

	// Get Public Key
	keyIDInt, err := strconv.Atoi(query.Get("key_id"))
	if err != nil {
		return nil, ErrCannotFoundPublicKey
	}
	callBackParam.KeyID = keyIDInt
	key, err := v.getKey(callBackParam.KeyID)
	if err != nil {
		return nil, err
	}

	// Verify Signature
	signatureBinary, err := base64.RawURLEncoding.DecodeString(callBackParam.Signature)
	if err != nil {
		return nil, ErrInvalidSignature
	}
	verified := ecdsa.VerifyASN1(key, msgHash, signatureBinary)
	if !verified {
		return nil, ErrInvalidSignature
	}

	return callBackParam, nil
}

func (v *Verifier) getKey(keyID int) (*ecdsa.PublicKey, error) {
	key, exist := v.cache.Get(strconv.Itoa(keyID))
	if !exist {
		err := v.loadKey()
		if err != nil {
			return nil, err
		}
		key, exist = v.cache.Get(strconv.Itoa(keyID))
		if !exist {
			return nil, ErrCannotFoundPublicKey
		}
	}

	ecdsaPublicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrCannotFoundPublicKey
	}

	return ecdsaPublicKey, nil
}

func (v *Verifier) loadKey() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	var keys adModResponse

	req, err := http.NewRequest(http.MethodGet, adModKeyServer, nil)
	if err != nil {
		return err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return err
	}

	for _, key := range keys.Keys {
		block, _ := pem.Decode([]byte(key.Pem))
		if block == nil {
			return ErrLoadPublicKey
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return ErrLoadPublicKey
		}

		ecdsaPublicKey, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			return ErrLoadPublicKey
		}

		v.cache.Set(strconv.Itoa(key.KeyId), ecdsaPublicKey, cache.DefaultExpiration)
	}

	return nil
}

// RefreshPublicKey Clean all cached Admob public keys
func (v *Verifier) RefreshPublicKey() {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.cache.Flush()
}

func hash(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	// compute the SHA256 hash
	return h.Sum(nil)
}
