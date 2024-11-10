package admob

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestAdmobVerify(t *testing.T) {
	verifier := NewVerifier()

	var tests = []struct {
		url           string
		rewardItem    string
		transactionID string
		customData    string
		userID        string
	}{
		{
			url:           "https://abc.exmaple.com/check?ad_network=5450213213286189855&ad_unit=1234567890&custom_data=test_cutom_data&reward_amount=1&reward_item=Test%20Reward&timestamp=1715323660511&transaction_id=123456789&user_id=test_user_id&signature=MEUCIHYdinf1Le0VSacI1cStAkAyBLas8eO8PKuZSt0ltOOAAiEAyoPvNIfTYqRrm7wCi_z-JjYEIiXvhLnqZu2fMpsRNIU&key_id=3335741209",
			rewardItem:    "Test Reward",
			transactionID: "123456789",
			customData:    "test_cutom_data",
			userID:        "test_user_id",
		},
		{
			url:           "https://abc.exmaple.com/check?ad_network=5450213213286189855&ad_unit=1234567890&custom_data=test%20cutom%20data&reward_amount=1&reward_item=Test%20Reward&timestamp=1715324253607&transaction_id=123456789&user_id=test%20user%20id&signature=MEUCIAyiWsF4RIaOIiWQpE0ABexqpkD4z0UpT0Hyvqyk-gyPAiEAnwrQVa3abQ7D2DvBBYyZmpFtmCyxZBtEEkizCqGf9oI&key_id=3335741209",
			rewardItem:    "Test Reward",
			transactionID: "123456789",
			customData:    "test cutom data",
			userID:        "test user id",
		},
	}

	err := verifier.loadKey()
	require.NoError(t, err)
	require.Greater(t, verifier.cache.ItemCount(), 0)

	for i := range tests {
		cbUrl, err := url.Parse(tests[i].url)
		callbackParam, err := verifier.Verify(cbUrl)
		require.NoError(t, err)

		require.Equal(t, tests[i].rewardItem, callbackParam.RewardItem)
		require.Equal(t, tests[i].transactionID, callbackParam.TransactionID)
		require.Equal(t, tests[i].customData, callbackParam.CustomData)
		require.Equal(t, tests[i].userID, callbackParam.UserID)
	}

	verifier.RefreshPublicKey()
	require.Equal(t, verifier.cache.ItemCount(), 0)
}
