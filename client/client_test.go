package client

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var client = Client{URL: "http://localhost:8080"}

func buildAccount() Account {
	return Account{
		Data: Data{
			Type:           "accounts",
			ID:             uuid.NewString(),
			OrganisationID: uuid.NewString(),
			Attributes: Attributes{
				Country: "GB",
			},
		},
	}
}

func TestCreate(t *testing.T) {
	account := buildAccount()
	resAccount, err := client.Create(account)

	assert.Equal(t, account.Data.ID, resAccount.Data.ID)
	assert.Equal(t, account.Data.OrganisationID, resAccount.Data.OrganisationID)
	assert.Equal(t, account.Data.Attributes.Country, resAccount.Data.Attributes.Country)
	assert.Nil(t, err)

	t.Cleanup(func() {
		require.Nil(t, client.Delete(account.Data.ID, account.Data.Version))
	})
}

func TestCreateWithDuplicateID(t *testing.T) {
	// set up account
	resAccount, err := client.Create(buildAccount())
	require.NotNil(t, resAccount)
	require.Nil(t, err)

	// test Create with duplicate ID
	dupAccount := buildAccount()
	dupAccount.Data.ID = resAccount.Data.ID
	resDupAccount, err := client.Create(dupAccount)
	assert.Nil(t, resDupAccount)
	assert.Equal(t,
		"409 Conflict: Account cannot be created as it violates a duplicate constraint",
		err.Error())

	t.Cleanup(func() {
		require.Nil(t, client.Delete(resAccount.Data.ID, resAccount.Data.Version))
	})
}

func TestCreateWithEmptyAccountFields(t *testing.T) {
	account := Account{}
	resAccount, err := client.Create(account)
	assert.Nil(t, resAccount)
	assert.Equal(t, "400 Bad Request: validation failure list:\n"+
		"validation failure list:\n"+
		"validation failure list:\n"+
		"country in body is required\n"+
		"id in body is required\n"+
		"organisation_id in body is required\n"+
		"type in body is required",
		err.Error())
}

func TestCreateWithNotValidRequiredAccountFields(t *testing.T) {
	account := Account{}
	account.Data.ID = "123abc"
	account.Data.OrganisationID = "123abc"
	account.Data.Type = "account"
	account.Data.Attributes.Country = "SVK"

	resAccount, err := client.Create(account)
	assert.Nil(t, resAccount)
	assert.Equal(t, "400 Bad Request: validation failure list:\n"+
		"validation failure list:\n"+
		"validation failure list:\n"+
		"country in body should match '^[A-Z]{2}$'\n"+
		"id in body must be of type uuid: \"123abc\"\n"+
		"organisation_id in body must be of type uuid: \"123abc\"\n"+
		"type in body should be one of [accounts]",
		err.Error())
}

func TestFetch(t *testing.T) {
	// set up account
	account, err := client.Create(buildAccount())
	require.NotNil(t, account)
	require.Nil(t, err)

	// test Fetch
	resAccount, err := client.Fetch(account.Data.ID)
	assert.Equal(t, account.Data.ID, resAccount.Data.ID)
	assert.Equal(t, account.Data.OrganisationID, resAccount.Data.OrganisationID)
	assert.Equal(t, account.Data.Attributes.Country, resAccount.Data.Attributes.Country)
	assert.Nil(t, err)

	t.Cleanup(func() {
		require.Nil(t, client.Delete(account.Data.ID, account.Data.Version))
	})
}

func TestFetchWitEmptyID(t *testing.T) {
	fetchWitIncorrectID(t, "", "Account.Data.ID is empty")
}

func TestFetchWitNotValidID(t *testing.T) {
	fetchWitIncorrectID(t, "123abc", "400 Bad Request: id is not a valid uuid")
}

func fetchWitIncorrectID(t *testing.T, accountID string, expectedErrorMessage string) {
	resAccount, err := client.Fetch(accountID)
	assert.Nil(t, resAccount)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestFetchWithNotExistingID(t *testing.T) {
	id := uuid.NewString()
	resAccount, err := client.Fetch(id)
	assert.Nil(t, resAccount)
	assert.Equal(t,
		fmt.Sprintf("404 Not Found: record %s does not exist", id),
		err.Error())
}

func TestDelete(t *testing.T) {
	// set up account
	account, err := client.Create(buildAccount())
	require.NotNil(t, account)
	require.Nil(t, err)

	// test Delete
	assert.Nil(t, client.Delete(account.Data.ID, account.Data.Version))

	// doublecheck the account is deleted
	_, err = client.Fetch(account.Data.ID)
	assert.Equal(t,
		fmt.Sprintf("404 Not Found: record %s does not exist", account.Data.ID),
		err.Error())
}

func TestDeleteWitEmptyID(t *testing.T) {
	err := client.Delete("", 0)
	assert.Equal(t, "404 Not Found", err.Error())
}

func TestDeleteWitNotValidID(t *testing.T) {
	err := client.Delete("123abc", 0)
	assert.Equal(t, "400 Bad Request: id is not a valid uuid", err.Error())
}

func TestDeleteWitNotExistingID(t *testing.T) {
	assert.Nil(t, client.Delete(uuid.NewString(), 0))
}

func TestDeleteWitWrongAccountVersion(t *testing.T) {
	// set up account
	account, err := client.Create(buildAccount())
	require.NotNil(t, account)
	require.Nil(t, err)

	// test Delete
	err = client.Delete(account.Data.ID, account.Data.Version+1)
	assert.Equal(t, "404 Not Found: invalid version", err.Error())

	t.Cleanup(func() {
		require.Nil(t, client.Delete(account.Data.ID, account.Data.Version))
	})
}

func TestWrongBaseURL(t *testing.T) {
	resAccount, err := (&Client{URL: "http://localhost:2468"}).Create(buildAccount())
	assert.Nil(t, resAccount)
	assert.Equal(t,
		`Post "http://localhost:2468/v1/organisation/accounts": dial tcp [::1]:2468: connect: connection refused`,
		err.Error())
}

func TestEmptyClientURL(t *testing.T) {
	resAccount, err := (&Client{}).Create(buildAccount())
	assert.Nil(t, resAccount)
	assert.Equal(t,
		`Post "/v1/organisation/accounts": unsupported protocol scheme ""`,
		err.Error())
}
