package admob

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func ExampleVerify() {
	// With Default Config
	verifier := NewVerifier()

	// With Custom Config
	verifier = NewVerifyWithConfig(&Config{
		PublicKeyCacheDuration: time.Hour * 10,
		HttpClient:             http.DefaultClient,
	})

	u, err := url.Parse("https://abc.exmaple.com/check?ad_network=5450213213286189855&ad_unit=1234567890&custom_data=test_cutom_data&reward_amount=1&reward_item=Test%20Reward&timestamp=1715323660511&transaction_id=123456789&user_id=test_user_id&signature=MEUCIHYdinf1Le0VSacI1cStAkAyBLas8eO8PKuZSt0ltOOAAiEAyoPvNIfTYqRrm7wCi_z-JjYEIiXvhLnqZu2fMpsRNIU&key_id=3335741209")
	if err != nil {
		panic(err)
	}

	callbackParam, err := verifier.Verify(u)
	if err != nil {
		panic(err)
	}

	fmt.Printf("AD Unit ID %s", callbackParam.AdUnit)

}
