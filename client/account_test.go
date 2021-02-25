package client

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var account = Account{
	Data: Data{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Version:        0,
		CreatedOn:      "2021-02-23T16:05:49.761Z",
		ModifiedOn:     "2021-02-23T16:05:49.761Z",
		Attributes: Attributes{
			Country:                 "GB",
			BaseCurrency:            "GBP",
			BankID:                  "400300",
			BankIDCode:              "GBDSC",
			BIC:                     "NWBKGB22",
			AccountNumber:           "41426819",
			IBAN:                    "GB11NWBK40030041426819",
			CustomerID:              "abc123",
			Name:                    [4]string{"Samantha Holder"},
			AlternativeNames:        [3]string{"Sam Holder"},
			AccountClassification:   "Personal",
			JointAccount:            false,
			AccountMatchingOptOut:   false,
			SecondaryIdentification: "A1B2C3D4",
			Switched:                false,
			Status:                  "confirmed",
		},
	},
	Links: Links{
		Self:  "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		First: "/v1/organisation/accounts?page%5Bnumber%5D=first",
		Last:  "/v1/organisation/accounts?page%5Bnumber%5D=last",
		Next:  "/v1/organisation/accounts?page%5Bnumber%5D=next",
		Prev:  "/v1/organisation/accounts?page%5Bnumber%5D=prev",
	},
}

func TestMarshalAccountToJSONStringAndBack(t *testing.T) {
	// build json string from an account
	json, err := account.json()
	require.NotEmpty(t, json)
	require.Nil(t, err)

	// build account from a json string
	jsonBytes := []byte(json)
	unmarshaledAccount := Account{}
	err = unmarshaledAccount.build(&jsonBytes)
	require.Nil(t, err)

	// compare the original with processed values
	assert.True(t, cmp.Equal(account, unmarshaledAccount))
}

func TestBuilAccountdWithNotValidJSONString(t *testing.T) {
	// create malformed json
	json, err := account.json()
	require.NotEmpty(t, json)
	require.Nil(t, err)
	malformedJSONBytes := []byte(json[1:])

	extractedAccount := Account{}
	err = extractedAccount.build(&malformedJSONBytes)
	assert.NotNil(t, err)
}
